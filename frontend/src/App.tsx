// frontend/src/App.tsx
import React, { useEffect, useState } from 'react';
import { fetchGeoIP } from './api';
import './App.css';

function App() {
  const [ip, setIp] = useState('217.150.32.5');
  const [geoData, setGeoData] = useState(null);
  const [error, setError] = useState('');

  const handleFetch = async () => {
    try {
      const data = await fetchGeoIP(ip);
      setGeoData(data);
      setError('');
    } catch (err) {
      setError(String(err));
      setGeoData(null);
    }
  };

  useEffect(() => {
    handleFetch();
  }, []);

  return (
      <div className="App">
        <h1>GeoIP Lookup</h1>
        <input
            value={ip}
            onChange={(e) => setIp(e.target.value)}
            placeholder="Enter IP address"
        />
        <button onClick={handleFetch}>Check IP</button>
        {error && <div className="error">{error}</div>}
        {geoData && <pre>{JSON.stringify(geoData, null, 2)}</pre>}
      </div>
  );
}

export default App;
