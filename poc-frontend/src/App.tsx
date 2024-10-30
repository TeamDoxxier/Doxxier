import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import TransferScreen from './TransferScreen';
import {Doxxier, FileThumbnail} from './types/Doxxier';

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
  const handleToggle = (id: string) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.Id === id ? { ...doxxier, isRevealed: !doxxier.IsRevealed } : doxxier
      )
    );
  };

  // Handle description input change
  const handleDescriptionChange = (id: string, value: string) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.Id === id ? { ...doxxier, description: value } : doxxier
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
          Url: isImage ? URL.createObjectURL(file) : '',
          Name: file.name,
          IsImage: isImage,
        };
      });

      setDoxxiers((prevDoxxiers) =>
        prevDoxxiers.map((doxxier, i) =>
          i === index
            ? {
                ...doxxier,
                thumbnails: [...doxxier.Thumbnails, ...fileThumbnails],
                thumbnailProgress: [...doxxier.ThumbnailProgress, ...Array(files.length).fill(0)], // Initialize progress for new files
              }
            : doxxier
        )
      );
    }
  };

  // Remove a specific file thumbnail
  const handleRemoveThumbnail = (doxxierId: string, thumbnailIndex: number) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.Id === doxxierId
          ? {
              ...doxxier,
              thumbnails: doxxier.Thumbnails.filter((_, index) => index !== thumbnailIndex),
              thumbnailProgress: doxxier.ThumbnailProgress.filter((_, index) => index !== thumbnailIndex),
            }
          : doxxier
      )
    );
  };

  // Create a new Doxxier container and auto-expand it
  const handleCreateNewDoxxier = () => {
      // Call the CreateDoxxier function from the Wasm module
      const doxxierData = window.CreateDoxxier();

      const goDoxxier = JSON.parse(doxxierData) as Doxxier;
      console.log('Wasm CreateDoxxier result:', doxxierData);

      const newDoxxier: Doxxier = {
        Id: goDoxxier.Id , // Use the ID from the Wasm module or fallback to newId
        Description: `Doxxier ${goDoxxier.Id}`,
        Original: '200MB', // Example values; use `doxxierData` if it returns these
        Packaged: '2KB',
        EstimatedDelivery: '3 days',
        IsTransferring: false,
        Progress: 0,
        Thumbnails: [],
        ThumbnailProgress: [],
        IsRevealed: true,
        Parts: [],
        CreatedAt: goDoxxier.CreatedAt,
      };
      setDoxxiers([...doxxiers, newDoxxier]);
  };

  // Handle 'Send' button click to start the transfer simulation
  const handleSend = (id: string) => {
    setDoxxiers((prevDoxxiers) =>
      prevDoxxiers.map((doxxier) =>
        doxxier.Id === id ? { ...doxxier, isTransferring: true } : doxxier
      )
    );

    // Simulate file transfer
    simulateTransfer(id);
  };

  // Simulate the transfer of files in a Doxxier container
  const simulateTransfer = (doxxierId: string) => {
    let currentFile = 0;

    const interval = setInterval(() => {
      setDoxxiers((prevDoxxiers) => {
        return prevDoxxiers.map((doxxier) => {
          if (doxxier.Id === doxxierId) {
            const newThumbnailProgress = [...doxxier.ThumbnailProgress];

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
  const handleRemoveDoxxier = (id: string) => {
    setDoxxiers((prevDoxxiers) => prevDoxxiers.filter((doxxier) => doxxier.Id !== id));
  };

  return (
    <div className="app-container">
      <div className="doxxier-list-container">
        {doxxiers.map((doxxier, index) => (
          <div key={doxxier.Id} className="doxxier-container">
            {doxxier.IsTransferring ? (
              <TransferScreen
                doxxierId={doxxier.Id}
                description={doxxier.Description}
                thumbnails={doxxier.Thumbnails}
                onClose={() => handleRemoveDoxxier(doxxier.Id)}
                thumbnailProgress={doxxier.ThumbnailProgress}
                overallProgress={doxxier.Progress}
              />
            ) : (
              <>
                <div className="header">
                  <span className="title clickable" onClick={() => handleToggle(doxxier.Id)}>
                    Doxxier {doxxier.Id} {doxxier.IsRevealed ? '▲' : '▼'}
                  </span>
                  <button className="remove-button" onClick={() => handleRemoveDoxxier(doxxier.Id)}>
                    X
                  </button>
                </div>

                <div className={`content-section ${doxxier.IsRevealed ? 'revealed' : 'hidden'}`}>
                  <input
                    type="text"
                    className="description-input"
                    placeholder="Enter description"
                    value={doxxier.Description}
                    onChange={(e) => handleDescriptionChange(doxxier.Id, e.target.value)}
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
                    {doxxier.Thumbnails.map((thumbnail, idx) => (
                      <div key={idx} className="thumbnail-container">
                        {thumbnail.IsImage ? (
                          <img src={thumbnail.Url} alt="thumbnail" className="thumbnail" />
                        ) : (
                          <div className="file-placeholder">{thumbnail.Name.split('.').pop()}</div>
                        )}
                        <button
                          className="thumbnail-remove-button"
                          onClick={() => handleRemoveThumbnail(doxxier.Id, idx)}
                        >
                          ×
                        </button>
                      </div>
                    ))}
                  </div>

                  <div className="status-section">
                    <p className="label">Original: {doxxier.Original}</p>
                    <p className="label">Packaged: {doxxier.Packaged}</p>
                    <p className="label">Estimated delivery: {doxxier.EstimatedDelivery}</p>
                  </div>

                  <div className="footer">
                    <button className="send-button" onClick={() => handleSend(doxxier.Id)}>
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
