// frontend/src/api.ts
export const fetchGeoIP = async (ip: string) => {
    const response = await fetch(`/api/geoip?ip=${ip}`);
    if (!response.ok) throw new Error('Error fetching GeoIP data');
    return await response.json();
};
