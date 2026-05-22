import { config } from '../config.js';
export function mapHttpStatusCode(status) {
    switch (status) {
        case 400: return 'BAD_REQUEST';
        case 401: return 'UNAUTHORIZED';
        case 403: return 'FORBIDDEN';
        case 404: return 'NOT_FOUND';
        case 429: return 'TOO_MANY_REQUESTS';
        case 503: return 'SERVICE_UNAVAILABLE';
        default: return 'INTERNAL_ERROR';
    }
}
/**
 * Hàm gọi Fake UIT Server và chuẩn hóa Response theo dạng MCP Envelope
 */
export async function callFakeServer(toolName, method, path, token, body, query) {
    let url = `${config.apiBaseUrl}${path}`;
    if (query && Object.keys(query).length > 0) {
        const params = new URLSearchParams();
        for (const [key, value] of Object.entries(query)) {
            if (value !== undefined)
                params.append(key, String(value));
        }
        url += `?${params.toString()}`;
    }
    const headers = {
        'Content-Type': 'application/json',
    };
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }
    const endpointStr = `${method.toUpperCase()} ${config.apiBaseUrl}${path}`;
    try {
        const response = await fetch(url, {
            method,
            headers,
            body: body ? JSON.stringify(body) : undefined,
        });
        let responseData = {};
        const text = await response.text();
        if (text) {
            try {
                responseData = JSON.parse(text);
            }
            catch (e) {
                responseData = { message: text };
            }
        }
        if (!response.ok) {
            return {
                ok: false,
                tool: toolName,
                endpoint: endpointStr,
                error: {
                    http_status: response.status,
                    code: mapHttpStatusCode(response.status),
                    message: responseData.message || responseData.error_code || 'Lỗi từ máy chủ'
                }
            };
        }
        return {
            ok: true,
            tool: toolName,
            endpoint: endpointStr,
            data: responseData.data || responseData,
            meta: { source: 'fake-uit-server' }
        };
    }
    catch (err) {
        return {
            ok: false,
            tool: toolName,
            endpoint: endpointStr,
            error: {
                http_status: 503,
                code: 'SERVICE_UNAVAILABLE',
                message: err.message || 'Không thể kết nối đến máy chủ UIT.'
            }
        };
    }
}
/**
 * Format payload thành chuẩn trả về cho MCP SDK
 */
export function formatMcpResponse(data) {
    return {
        content: [{
                type: "text",
                text: JSON.stringify(data, null, 2)
            }]
    };
}
