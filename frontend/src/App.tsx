import React, { useState } from 'react';
import './App.css';

function App() {
    const [ips, setIps] = useState('');
    const [file, setFile] = useState<File | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();

        const formData = new FormData();
        if (ips) formData.append('ips', ips);
        if (file) formData.append('file', file);

        const response = await fetch('/check_ips', {
            method: 'POST',
            body: formData,
        });

        if (response.ok) {
            const blob = await response.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'geoip_results.csv';
            a.click();
        } else {
            alert('Error: ' + (await response.text()));
        }
    };

    return (
        <div className="App">
            <h1>GeoIP Batch Lookup</h1>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Paste up to 100 IPs (comma-separated):</label><br />
                    <textarea
                        value={ips}
                        onChange={(e) => setIps(e.target.value)}
                        rows={4}
                        cols={60}
                        placeholder="192.168.1.1, 8.8.8.8, ..."
                    />
                </div>

                <div>
                    <label>Or upload a .txt or .csv file:</label><br />
                    <input
                        type="file"
                        accept=".txt,.csv"
                        onChange={(e) => {
                            const files = e.target.files;
                            if (files && files.length > 0) setFile(files[0]);
                        }}
                    />
                </div>

                <button type="submit">Submit</button>
            </form>
        </div>
    );
}

export default App;
