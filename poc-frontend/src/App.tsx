import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import TransferScreen from './TransferScreen';

interface FileThumbnail {
  url: string;
  name: string;
  isImage: boolean;
}

interface Doxxier {
  id: number;
  original: string;
  packaged: string;
  estimatedDelivery: string;
  isRevealed: boolean;
  isTransferring: boolean;
  progress: number; // Overall transfer progress
  description: string;
  thumbnails: FileThumbnail[];
  thumbnailProgress: number[]; // Individual file progress
}

const App: React.FC = () => {
  const [doxxiers, setDoxxiers] = useState<Doxxier[]>([]);
  const [connectionStatus, setConnectionStatus] = useState('Connecting to cMixx'); // Status text
  const [isConnected, setIsConnected] = useState(false); // Connection state
  const fileInputRefs = useRef<HTMLInputElement[]>([]);

  useEffect(() => {
    // Animate the "Connecting to cMixx..." text
    let animationInterval: NodeJS.Timeout;
    let statusIndex = 0;
    const statusMessages = ['Connecting to cMixx', 'Connecting to cMixx.', 'Connecting to cMixx..', 'Connecting to cMixx...'];
    
    if (!isConnected) {
      animationInterval = setInterval(() => {
        statusIndex = (statusIndex + 1) % statusMessages.length;
        setConnectionStatus(statusMessages[statusIndex]);
      }, 500); // Update the animation every 500ms

      // Switch to "Connected" after 5 seconds
      setTimeout(() => {
        clearInterval(animationInterval);
        setConnectionStatus('Connected to cMixx');
        setIsConnected(true);
      }, 5000);
    }

    return () => {
      if (animationInterval) clearInterval(animationInterval);
    };
  }, [isConnected]);

  // Toggle the visibility of a Doxxier's content
  const handleToggle = (id: number) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.id === id ? { ...doxxier, isRevealed: !doxxier.isRevealed } : doxxier
      )
    );
  };

  // Handle description input change
  const handleDescriptionChange = (id: number, value: string) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.id === id ? { ...doxxier, description: value } : doxxier
      )
    );
  };

  // Handle adding a file when "Add Files" is clicked
  const handleAddClick = (index: number) => {
    const fileInput = fileInputRefs.current[index];
    if (fileInput) {
      fileInput.click();
    }
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>, index: number) => {
    if (event.target.files) {
      const files = Array.from(event.target.files);
      const fileThumbnails: FileThumbnail[] = files.map((file) => {
        const isImage = file.type.startsWith('image/');
        return {
          url: isImage ? URL.createObjectURL(file) : '',
          name: file.name,
          isImage,
        };
      });

      setDoxxiers((prevDoxxiers) =>
        prevDoxxiers.map((doxxier, i) =>
          i === index
            ? {
                ...doxxier,
                thumbnails: [...doxxier.thumbnails, ...fileThumbnails],
                thumbnailProgress: [...doxxier.thumbnailProgress, ...Array(files.length).fill(0)], // Initialize progress for new files
              }
            : doxxier
        )
      );
    }
  };

  // Remove a specific file thumbnail
  const handleRemoveThumbnail = (doxxierId: number, thumbnailIndex: number) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.id === doxxierId
          ? {
              ...doxxier,
              thumbnails: doxxier.thumbnails.filter((_, index) => index !== thumbnailIndex),
              thumbnailProgress: doxxier.thumbnailProgress.filter((_, index) => index !== thumbnailIndex),
            }
          : doxxier
      )
    );
  };

  // Create a new Doxxier container and auto-expand it
  const handleCreateNewDoxxier = () => {
    const newId = doxxiers.length + 1;
    setDoxxiers([
      ...doxxiers,
      {
        id: newId,
        original: '200MB',
        packaged: '2KB',
        estimatedDelivery: '3 days',
        isRevealed: true,
        isTransferring: false,
        progress: 0,
        description: '',
        thumbnails: [],
        thumbnailProgress: [], // Initialize with empty progress
      },
    ]);
  };

  // Handle 'Send' button click to start the transfer simulation
  const handleSend = (id: number) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.id === id ? { ...doxxier, isTransferring: true } : doxxier
      )
    );

    // Simulate file transfer
    simulateTransfer(id);
  };

  // Simulate the transfer of files in a Doxxier container
  const simulateTransfer = (doxxierId: number) => {
    let currentFile = 0;

    const interval = setInterval(() => {
      setDoxxiers((prevDoxxiers) => {
        return prevDoxxiers.map((doxxier) => {
          if (doxxier.id === doxxierId) {
            const newThumbnailProgress = [...doxxier.thumbnailProgress];

            // Increment the progress of the current file
            if (newThumbnailProgress[currentFile] < 100) {
              newThumbnailProgress[currentFile] += 10; // Increment by 10%
            } else {
              // Move to the next file if the current one is done
              currentFile += 1;
            }

            // Calculate overall progress
            const totalProgress = newThumbnailProgress.reduce((acc, val) => acc + val, 0);
            const overallProgress = (totalProgress / (newThumbnailProgress.length * 100)) * 100;

            // Stop the interval if all files are processed
            if (currentFile >= newThumbnailProgress.length) {
              clearInterval(interval);
            }

            return {
              ...doxxier,
              thumbnailProgress: newThumbnailProgress,
              progress: overallProgress,
            };
          }
          return doxxier;
        });
      });
    }, 1000); // Update every second
  };

  // Remove a specific Doxxier container
  const handleRemoveDoxxier = (id: number) => {
    setDoxxiers((prevDoxxiers) => prevDoxxiers.filter((doxxier) => doxxier.id !== id));
  };

  return (
    <div className="app-container">
      <div className="doxxier-list-container">
        {doxxiers.map((doxxier, index) => (
          <div key={doxxier.id} className="doxxier-container">
            {doxxier.isTransferring ? (
              <TransferScreen
                doxxierId={doxxier.id}
                description={doxxier.description}
                thumbnails={doxxier.thumbnails}
                onClose={() => handleRemoveDoxxier(doxxier.id)}
                thumbnailProgress={doxxier.thumbnailProgress}
                overallProgress={doxxier.progress}
              />
            ) : (
              <>
                <div className="header">
                  <span className="title clickable" onClick={() => handleToggle(doxxier.id)}>
                    Doxxier {doxxier.id} {doxxier.isRevealed ? '▲' : '▼'}
                  </span>
                  <button className="remove-button" onClick={() => handleRemoveDoxxier(doxxier.id)}>
                    X
                  </button>
                </div>

                <div className={`content-section ${doxxier.isRevealed ? 'revealed' : 'hidden'}`}>
                  <input
                    type="text"
                    className="description-input"
                    placeholder="Enter description"
                    value={doxxier.description}
                    onChange={(e) => handleDescriptionChange(doxxier.id, e.target.value)}
                  />

                  <div className="add-section">
                    <button className="add-button" onClick={() => handleAddClick(index)}>
                      Add Files
                    </button>
                    <input
                      type="file"
                      ref={(el) => (fileInputRefs.current[index] = el!)}
                      onChange={(e) => handleFileChange(e, index)}
                      style={{ display: 'none' }}
                      multiple
                      accept="image/*,.doc,.pdf,.txt"
                    />
                    {/* Render thumbnails and placeholders */}
                    {doxxier.thumbnails.map((thumbnail, idx) => (
                      <div key={idx} className="thumbnail-container">
                        {thumbnail.isImage ? (
                          <img src={thumbnail.url} alt="thumbnail" className="thumbnail" />
                        ) : (
                          <div className="file-placeholder">{thumbnail.name.split('.').pop()}</div>
                        )}
                        <button
                          className="thumbnail-remove-button"
                          onClick={() => handleRemoveThumbnail(doxxier.id, idx)}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>

                  <div className="status-section">
                    <p className="label">Original: {doxxier.original}</p>
                    <p className="label">Packaged: {doxxier.packaged}</p>
                    <p className="label">Estimated delivery: {doxxier.estimatedDelivery}</p>
                  </div>

                  <div className="footer">
                    <button className="send-button" onClick={() => handleSend(doxxier.id)}>
                      Send
                    </button>
                    <button className="save-button">Save</button>
                  </div>
                </div>
              </>
            )}
          </div>
        ))}

        {/* New Doxxier Button */}
        <div className="new-doxxier-container">
          <button className="new-doxxier-button" onClick={handleCreateNewDoxxier}>
            New Doxxier
          </button>
        </div>

        {/* Separator and connection status at the bottom of the container */}
        <div className="separator"></div>
        <div className="connection-status-container">
          <span className={`connection-status ${isConnected ? 'connected' : 'connecting'}`}>
            {connectionStatus}
          </span>
        </div>
      </div>
    </div>
  );
};

export default App;
