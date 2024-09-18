// index.tsx
import React from 'react';
import ReactDOM from 'react-dom/client'; // Import createRoot from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'; // Make sure to wrap App in BrowserRouter here
import App from './App';
import './index.css';

const root = ReactDOM.createRoot(document.getElementById('root') as HTMLElement);
root.render(
  <React.StrictMode>
    <BrowserRouter> {/* Wrap the entire app in BrowserRouter */}
      <App />
    </BrowserRouter>
  </React.StrictMode>
);
