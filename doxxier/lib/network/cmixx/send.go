////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file                                                               //
////////////////////////////////////////////////////////////////////////////////

package cmixx

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/v4/cmix"
	"gitlab.com/elixxir/client/v4/cmix/message"
	"gitlab.com/elixxir/client/v4/cmix/rounds"
	"gitlab.com/elixxir/crypto/dm"
	"gitlab.com/elixxir/crypto/fastRNG"
	"gitlab.com/elixxir/crypto/nike"
	"gitlab.com/elixxir/crypto/nike/ecdh"
	"gitlab.com/elixxir/primitives/format"
	"gitlab.com/xx_network/primitives/id"
	"gitlab.com/xx_network/primitives/id/ephemeral"
	"gitlab.com/xx_network/primitives/netTime"
	"golang.org/x/crypto/blake2b"
	"google.golang.org/protobuf/proto"
)

const (
	// Versions for various message types
	textVersion       = 0
	reactionVersion   = 0
	invitationVersion = 0
	silentVersion     = 0
	deleteVersion     = 0

	// SendMessageTag is the base tag used when generating a debug tag for
	// sending a message.
	SendMessageTag = "Message"

	// SendReplyTag is the base tag used when generating a debug tag for
	// sending a reply.
	SendReplyTag = "Reply"

	// SendReactionTag is the base tag used when generating a debug tag for
	// sending a reaction.
	SendReactionTag = "Reaction"

	// SendSilentTag is the base tag used when generating a debug tag for
	// sending a silent message.
	SendSilentTag = "Silent"

	// SendInviteTag is the base tag used when generating a debug tag for
	// sending an invitation.
	SendInviteTag = "Invite"

	// DeleteMessageTag is the base tag used when generating a debug tag for
	// delete message.
	DeleteMessageTag = "Delete"

	directMessageDebugTag = "dm"
	// The size of the nonce used in the message ID.
	messageNonceSize = 4
)

var (
	emptyPubKeyBytes = make([]byte, ed25519.PublicKeySize)
	emptyPubKey      = ed25519.PublicKey(emptyPubKeyBytes)
)

// Send is used to send a raw direct message to a DM partner. In general, it
// should be wrapped in a function that defines the wire protocol.
//
// If the final message, before being sent over the wire, is too long, this will
// return an error. Due to the underlying encoding using compression, it is not
// possible to define the largest payload that can be sent, but it will always
// be possible to send a payload of 802 bytes at minimum.
// DeleteMessage is used to send a formatted message to another user.
func (c *CMixxTransport) send(messageType MessageType, msg []byte) (rounds.Round, ephemeral.Id, error) {
	if c.partnerToken == 0 {
		return rounds.Round{},
			ephemeral.Id{},
			errors.Errorf("invalid dmToken value: %d", c.partnerToken)

	}

	if c.partnerPublicKey == nil ||
		c.partnerPublicKey.Equal(emptyPubKey) {
		return rounds.Round{},
			ephemeral.Id{},
			errors.Errorf("invalid public key value: %v",
				c.partnerPublicKey)
	}

	// if dc.myToken == partnerToken &&
	// 	!dc.me.PubKey.Equal(partnerEdwardsPubKey) {
	// 	return cryptoMessage.ID{}, rounds.Round{},
	// 		ephemeral.Id{},
	// 		errors.Errorf("can only use myToken on self send: "+
	// 			"myToken: %d, myKey: %v, partnerKey: %v, partnerToken: %d",
	// 			dc.myToken, dc.me.PubKey, partnerEdwardsPubKey, partnerToken)
	// }

	partnerPubKey := ecdh.Edwards2EcdhNikePublicKey(c.partnerPublicKey)

	partnerID := deriveReceptionID(partnerPubKey.Bytes(), c.partnerToken)

	sihTag := dm.MakeSenderSihTag(c.partnerPublicKey, c.directId.Privkey)
	mt := messageType.Marshal()
	service := message.CompressedService{
		Identifier: c.partnerPublicKey,
		Tags:       []string{sihTag},
		Metadata:   mt[:],
	}

	// Note: We log sends on exit, and append what happened to the message
	// this cuts down on clutter in the log.
	params := cmix.GetDefaultCMIXParams()
	tag := makeDebugTag(c.partnerPublicKey, []byte(msg), SendReplyTag)

	params = params.SetDebugTag(tag)
	sendPrint := fmt.Sprintf("[DM][%s] Sending from %s to %s type %s at %s",
		params.DebugTag, base64.StdEncoding.EncodeToString(c.directId.PubKey),
		partnerID, messageType,
		netTime.Now())
	defer func() { jww.INFO.Println(sendPrint) }()

	rng := c.cMixxNet.GetRng().GetStream()
	defer rng.Close()

	// Generate random nonce to be used for message ID
	// generation. This makes it so two identical messages sent on
	// the same round have different message IDs.
	msgNonce := make([]byte, messageNonceSize)
	n, err := rng.Read(msgNonce)
	if err != nil {
		sendPrint += fmt.Sprintf(", failed to generate nonce: %+v", err)
		return rounds.Round{},
			ephemeral.Id{},
			errors.Errorf("Failed to generate nonce: %+v", err)
	} else if n != messageNonceSize {
		sendPrint += fmt.Sprintf(
			", got %d bytes for %d-byte nonce", n, messageNonceSize)
		return rounds.Round{},
			ephemeral.Id{},
			errors.Errorf(
				"Generated %d bytes for %d-byte nonce", n,
				messageNonceSize)
	}

	directMessage := &DirectMessage{
		DmToken:        c.directId.GetDMToken(),
		PayloadType:    uint32(messageType),
		Payload:        msg,
		Nonce:          msgNonce,
		LocalTimestamp: netTime.Now().UnixNano(),
	}

	if params.DebugTag == cmix.DefaultDebugTag {
		params.DebugTag = directMessageDebugTag
	}

	sendPrint += fmt.Sprintf(", pending send %s", netTime.Now())
	//uuid, err := dc.st.DenotePendingSend(partnerEdwardsPubKey,
	//dc.me.PubKey, partnerToken, messageType, directMessage)
	// if err != nil {
	// 	sendPrint += fmt.Sprintf(", pending send failed %s",
	// 		err.Error())
	// 	errDenote := dc.st.FailedSend(uuid)
	// 	if errDenote != nil {
	// 		sendPrint += fmt.Sprintf(
	// 			", failed to denote failed dm send: %s",
	// 			errDenote.Error())
	// 	}
	// 	return cryptoMessage.ID{}, rounds.Round{},
	// 		ephemeral.Id{}, err
	// }

	rndID, ephIDs, err := send(c.cMixxNet.GetCmix(), c.selfReceptionId,
		partnerID, partnerPubKey, c.privateKey, service,
		c.partnerToken, directMessage, params, c.cMixxNet.GetRng())
	if err != nil {
		sendPrint += fmt.Sprintf(", err on send: %+v", err)
		// errDenote := dc.st.FailedSend(uuid)
		// if errDenote != nil {
		// 	sendPrint += fmt.Sprintf(
		// 		", failed to denote failed dm send: %s",
		// 		errDenote.Error())
		// }
		// return cryptoMessage.ID{}, rounds.Round{},
		// 	ephemeral.Id{}, err
	}

	// Now that we have a round ID, derive the msgID
	// FIXME: cryptoMesage.DeriveDirectMessageID should take a round ID,
	// and the callee shouldn't have been modifying the data we sent.
	directMessage.RoundId = uint64(rndID.ID)
	jww.INFO.Printf("[DM] DeriveDirectMessage(%s...) Send", partnerID)

	//err = dc.st.Sent(uuid, msgID, rndID)
	// if err != nil {
	// 	sendPrint += fmt.Sprintf(", dm send denote failed: %s ",
	// 		err.Error())
	// }
	return rndID, ephIDs[1], err

}

// DeriveReceptionID returns a reception ID for direct messages sent
// to the user. It generates this ID by hashing the public key and
// an arbitrary idToken together. The ID type is set to "User".
func DeriveReceptionID(publicKey ed25519.PublicKey, idToken uint32) *id.ID {
	nikePubKey := ecdh.Edwards2EcdhNikePublicKey(publicKey)
	return deriveReceptionID(nikePubKey.Bytes(), idToken)
}

func deriveReceptionID(keyBytes []byte, idToken uint32) *id.ID {
	h, err := blake2b.New256(nil)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}
	h.Write(keyBytes)
	tokenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(tokenBytes, idToken)
	h.Write(tokenBytes)
	idBytes := h.Sum(nil)
	idBytes = append(idBytes, byte(id.User))
	receptionID, err := id.Unmarshal(idBytes)
	if err != nil {
		jww.FATAL.Panicf("%+v", err)
	}
	return receptionID
}

func send(net Cmix, myID *id.ID, partnerID *id.ID,
	partnerPubKey nike.PublicKey, myPrivateKey nike.PrivateKey,
	service cmix.Service, partnerToken uint32,
	msg *DirectMessage, params cmix.CMIXParams,
	rngGenerator *fastRNG.StreamGenerator) (rounds.Round,
	[]ephemeral.Id, error) {

	// Send to Partner
	assemble := func(rid id.Round) ([]cmix.TargetedCmixMessage, error) {
		rng := rngGenerator.GetStream()
		defer rng.Close()

		// Copy msg to dmMsg, which leaves the original
		// message data alone for resend purposes.
		// (deep copy isn't necessary because we only
		// change the rid)
		dmMsg := *msg

		// SEND
		dmMsg.RoundId = uint64(rid)

		// Serialize the message
		dmSerial, err := proto.Marshal(&dmMsg)
		if err != nil {
			return nil, err
		}

		payloadLen := calcDMPayloadLen(net)

		ciphertext := dm.Cipher.Encrypt(dmSerial, myPrivateKey,
			partnerPubKey, rng, payloadLen)

		fpBytes, encryptedPayload, mac, err := createCMIXFields(
			ciphertext, payloadLen, rng)
		if err != nil {
			return nil, err
		}

		fp := format.NewFingerprint(fpBytes)

		sendMsg := cmix.TargetedCmixMessage{
			Recipient:   partnerID,
			Payload:     encryptedPayload,
			Fingerprint: fp,
			Service:     service,
			Mac:         mac,
		}

		return []cmix.TargetedCmixMessage{sendMsg}, nil
	}
	return net.SendManyWithAssembler([]*id.ID{partnerID, myID}, assemble, params)
}

// makeDebugTag is a debug helper that creates non-unique msg identifier.
//
// This is set as the debug tag on messages and enables some level of tracing a
// message (if its contents/chan/type are unique).
func makeDebugTag(id ed25519.PublicKey,
	msg []byte, baseTag string) string {

	h, _ := blake2b.New256(nil)
	h.Write(msg)
	h.Write(id)

	tripCode := base64.RawStdEncoding.EncodeToString(h.Sum(nil))[:12]
	return fmt.Sprintf("%s-%s", baseTag, tripCode)
}

func calcDMPayloadLen(net Cmix) int {
	// As we don't use the mac or fp fields, we can extend
	// our payload size
	// (-2 to eliminate the first byte of mac and fp)
	return net.GetMaxMessageLength() +
		format.MacLen + format.KeyFPLen - 2

}

// Helper function that splits up the ciphertext into the appropriate cmix
// packet subfields
func createCMIXFields(ciphertext []byte, payloadSize int,
	rng io.Reader) (fpBytes, encryptedPayload, mac []byte, err error) {

	fpBytes = make([]byte, format.KeyFPLen)
	mac = make([]byte, format.MacLen)
	encryptedPayload = make([]byte, payloadSize-
		len(fpBytes)-len(mac)+2)

	// The first byte of mac and fp are random
	prefixBytes := make([]byte, 2)
	n, err := rng.Read(prefixBytes)
	if err != nil || n != len(prefixBytes) {
		err = fmt.Errorf("rng read failure: %+v", err)
		return nil, nil, nil, err
	}
	// Note: the first bit must be 0 for these...
	fpBytes[0] = 0x7F & prefixBytes[0]
	mac[0] = 0x7F & prefixBytes[1]

	// ciphertext[0:FPLen-1] == fp[1:FPLen]
	start := 0
	end := format.KeyFPLen - 1
	copy(fpBytes[1:format.KeyFPLen], ciphertext[start:end])
	// ciphertext[FPLen-1:FPLen+MacLen-2] == mac[1:MacLen]
	start = end
	end = start + format.MacLen - 1
	copy(mac[1:format.MacLen], ciphertext[start:end])
	// ciphertext[FPLen+MacLen-2:] == encryptedPayload
	start = end
	end = start + len(encryptedPayload)
	copy(encryptedPayload, ciphertext[start:end])

	// Fill anything left w/ random data
	numLeft := end - start - len(encryptedPayload)
	if numLeft > 0 {
		jww.WARN.Printf("[DM] small ciphertext, added %d bytes",
			numLeft)
		n, err := rng.Read(encryptedPayload[end-start:])
		if err != nil || n != numLeft {
			err = fmt.Errorf("rng read failure: %+v", err)
			return nil, nil, nil, err
		}
	}

	return fpBytes, encryptedPayload, mac, nil
}

func createRandomService(rng io.Reader) message.Service {
	// NOTE: 64 is entirely arbitrary, 33 bytes are used for the ID
	// and the rest will be base64'd into a string for the tag.
	data := make([]byte, 64)
	n, err := rng.Read(data)
	if err != nil {
		jww.FATAL.Panicf("rng failure: %+v", err)
	}
	if n != len(data) {
		jww.FATAL.Panicf("rng read failure, short read: %d < %d", n,
			len(data))
	}
	return message.Service{
		Identifier: data[:33],
		Tag:        base64.RawStdEncoding.EncodeToString(data[33:]),
	}
}
