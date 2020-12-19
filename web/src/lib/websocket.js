export default function(serverBaseUrl) {
    const  ws = new WebSocket(`${serverBaseUrl}/notifications`);
    return ws;
}