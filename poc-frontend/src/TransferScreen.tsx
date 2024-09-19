import React from 'react';
import './TransferScreen.css';
import { FileThumbnail } from './types/Doxxier';

interface TransferScreenProps {
  doxxierId: string;
  description: string;
  thumbnails: FileThumbnail[];
  thumbnailProgress: number[]; // Individual file progress
  overallProgress: number; // Overall progress
  onClose: () => void;
}

const TransferScreen: React.FC<TransferScreenProps> = ({
  doxxierId,
  description,
  thumbnails,
  thumbnailProgress,
  overallProgress,
  onClose,
}) => {
  return (
    <div className="transfer-screen-container">
      <button className="transfer-close-button" onClick={onClose}>
        ×
      </button>
      <div className="description-area">
        <p className="description-text">{description}</p>
      </div>
      <div className="thumbnail-section">
        {/* Display thumbnails with simulated progress */}
        {thumbnails.map((thumbnail, idx) => (
          <div key={idx} className="thumbnail-container" title={thumbnail.Name}>
            {thumbnail.IsImage ? (
              <div className="thumbnail-wrapper">
                <img src={thumbnail.Url} alt="thumbnail" className="thumbnail" />
                {/* Overlay to show the progress */}
                <div
                  className="thumbnail-progress-overlay"
                  style={{ transform: `translateX(${thumbnailProgress[idx]}%)` }} // Move overlay from left to right
                ></div>
              </div>
            ) : (
              <div className="thumbnail-wrapper">
                <div className="file-placeholder">{thumbnail.Name.split('.').pop()}</div>
                {/* Overlay to show the progress */}
                <div
                  className="thumbnail-progress-overlay"
                  style={{ transform: `translateX(${thumbnailProgress[idx]}%)` }} // Move overlay from left to right
                ></div>
              </div>
            )}
          </div>
        ))}
      </div>
      <div className="progress-indicator">
        <p>Transferring {Math.round(overallProgress)}%</p>
      </div>
    </div>
  );
};

export default TransferScreen;
