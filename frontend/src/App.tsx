import React, { useState } from 'react';
import './App.css';

type Result = {
    ip: string;
    country: string;
};

function App() {
    const [ips, setIps] = useState('');
    const [file, setFile] = useState<File | null>(null);
    const [results, setResults] = useState<Result[]>([]);
    const [error, setError] = useState('');
    const [showAll, setShowAll] = useState(false);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setResults([]);
        setError('');

        const formData = new FormData();
        if (ips) formData.append('ips', ips);
        if (file) formData.append('file', file);

        try {
            const response = await fetch('/check_ips', {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                const text = await response.text();
                throw new Error(text);
            }

            const text = await response.text();
            const rows = text.trim().split('\n').slice(1); // skip header
            const parsed: Result[] = rows.map(row => {
                const [ip, country] = row.split(',');
                return { ip, country };
            });
            setResults(parsed);
        } catch (err: any) {
            setError(err.message);
        }
    };

    const handleDownload = async () => {
        const formData = new FormData();
        if (ips) formData.append('ips', ips);
        if (file) formData.append('file', file);
        formData.append('download', '1');

        const res = await fetch('/check_ips', {
            method: 'POST',
            body: formData,
        });

        if (res.ok) {
            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'geoip_results.csv';
            a.click();
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
                {results.length > 0 && (
                    <>
                        <button type="button" onClick={handleDownload}>Download CSV</button>
                    </>
                )}
            </form>

            {error && <p className="error">{error}</p>}

            {results.length > 0 && (
                <div style={{ marginTop: '2rem' }}>
                    <h2>Results ({showAll ? results.length : 50} shown)</h2>
                    <table border={1} style={{ margin: 'auto', borderCollapse: 'collapse' }}>
                        <thead>
                        <tr>
                            <th>IP</th>
                            <th>Country Code</th>
                        </tr>
                        </thead>
                        <tbody>
                        {(showAll ? results : results.slice(0, 50)).map((res, i) => (
                            <tr key={i}>
                                <td>{res.ip}</td>
                                <td>{res.country}</td>
                            </tr>
                        ))}
                        </tbody>
                    </table>
                    {results.length > 50 && (
                        <button onClick={() => setShowAll(prev => !prev)}>
                            {showAll ? 'Show First 50' : 'Show All'}
                        </button>
                    )}
                </div>
            )}
        </div>
    );
}

export default App;
