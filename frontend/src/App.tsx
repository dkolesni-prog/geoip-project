import React, {useState} from 'react';
import './App.css';

type Result = {
    ip: string;
    country: string;
    countryName: string;
};

function isValidIP(ip: string): boolean {
    const ipv4Regex = /^(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d{2}|[1-9]?\d)){3}$/;
    const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$/;
    return ipv4Regex.test(ip) || ipv6Regex.test(ip);
}

function App() {
    const [ips, setIps] = useState('');
    const [invalidIPs, setInvalidIPs] = useState<string[]>([]);
    const [file, setFile] = useState<File | null>(null);
    const [results, setResults] = useState<Result[]>([]);
    const [error, setError] = useState('');
    const [showAll, setShowAll] = useState(false);

    const validateIps = (input: string) => {
        const list = input.split(',').map(ip => ip.trim()).filter(Boolean);
        const invalid = list.filter(ip => !isValidIP(ip));
        setInvalidIPs(invalid);
    };

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

            // ✅ We now assume it’s plain CSV text for preview, not a forced download
            const text = await response.text();
            const rows = text.trim().split('\n').slice(1); // skip header
            const parsed: Result[] = rows.map(row => {
                const [ip, country, countryName] = row.split(',');
                return {ip, country, countryName};
            });
            setResults(parsed);
            setShowAll(false);
        } catch (err: any) {
            setError(err.message);
        }
    };


    const handleDownload = async (format: 'csv' | 'mmdb') => {
        const formData = new FormData();
        if (ips) formData.append('ips', ips);
        if (file) formData.append('file', file);
        formData.append('export', format);

        const res = await fetch('/check_ips', {
            method: 'POST',
            body: formData,
        });

        if (res.ok) {
            const blob = await res.blob();
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `geoip_results.${format}`;
            a.click();
        } else {
            alert('Failed to download ' + format.toUpperCase());
        }
    };


    return (
        <div className="App">
            <h1>GeoIP Batch Lookup</h1>
            <form onSubmit={handleSubmit}>
                <div>
                    <label>Paste up to 100 IPs (comma-separated):</label><br/>
                    <textarea
                        value={ips}
                        onChange={(e) => {
                            const val = e.target.value;
                            setIps(val);
                            validateIps(val);
                        }}
                        rows={4}
                        cols={60}
                        placeholder="192.168.1.1, 8.8.8.8, ..."
                    />
                </div>

                <div>
                    <label>Or upload a .txt (comma or newline separated) or .json file (array of IPs):</label><br/>
                    <input
                        type="file"
                        accept=".txt,.json"
                        onChange={(e) => {
                            const files = e.target.files;
                            if (files && files.length > 0) setFile(files[0]);
                        }}
                    />
                </div>

                <button type="submit" disabled={invalidIPs.length > 0}>Check IPs</button>

                {results.length > 0 && (
                    <div style={{marginTop: '1rem'}}>
                        <button type="button" onClick={() => handleDownload('csv')}>Download CSV</button>
                        <button type="button" onClick={() => handleDownload('mmdb')}>Download MMDB</button>
                    </div>
                )}
            </form>


            {invalidIPs.length > 0 && (
                <div className="error" style={{marginTop: '1rem', textAlign: 'left'}}>
                    <p>⚠️ Invalid IPs detected:</p>
                    <ul>
                        {invalidIPs.map((ip, idx) => (
                            <li key={idx} style={{color: 'red'}}>{ip}</li>
                        ))}
                    </ul>
                </div>
            )}

            {error && <p className="error">{error}</p>}

            {results.length > 0 && (
                <div style={{marginTop: '2rem'}}>
                    <h2>Results ({showAll ? results.length : Math.min(results.length, 50)} shown)</h2>
                    <table border={1} style={{margin: 'auto', borderCollapse: 'collapse'}}>
                        <thead>
                        <tr>
                            <th>IP</th>
                            <th>Country Code</th>
                            <th>Country Name</th>
                            {/* ✅ New column */}
                        </tr>
                        </thead>
                        <tbody>
                        {(showAll ? results : results.slice(0, 50)).map((res, i) => (
                            <tr key={i}>
                                <td>{res.ip}</td>
                                <td>{res.country}</td>
                                <td>{res.countryName}</td>
                                {/* ✅ New column */}
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
