// Types based on proto definitions
export enum Protocol {
    PROTOCOL_UNSPECIFIED = 0,
    PROTOCOL_HTTP = 1,
    PROTOCOL_HTTP_STREAM = 2,
    PROTOCOL_GRPC = 3,
    PROTOCOL_GRPC_STREAM = 4,
}

export enum DeviceStatus {
    DEVICE_STATUS_UNSPECIFIED = "DEVICE_STATUS_UNSPECIFIED",
    DEVICE_STATUS_HEALTHY = "DEVICE_STATUS_HEALTHY",
    DEVICE_STATUS_DEGRADED = "DEVICE_STATUS_DEGRADED",
    DEVICE_STATUS_ERROR = "DEVICE_STATUS_ERROR",
    DEVICE_STATUS_MAINTENANCE = "DEVICE_STATUS_MAINTENANCE",
    DEVICE_STATUS_BOOTING = "DEVICE_STATUS_BOOTING",
    DEVICE_STATUS_OFFLINE = "DEVICE_STATUS_OFFLINE",
}

export interface Device {
    id: string;
    device_id: string;
    alias: string;
    host: string;
    port: number;
    port_gateway: number;
    architecture: string;
    os: string;
    supported_protocols: Protocol[];
    created_at?: string;
    updated_at?: string;
}

export interface Diagnostics {
    hardware_version: string;
    software_version: string;
    firmware_version: string;
    cpu_usage: number;
    memory_usage: number;
    device_status: DeviceStatus;
    checksum: string;
}

export interface DiagnosticsResponse {
    device?: Device;
    diagnostics?: Diagnostics;
    updated_at?: string;
}

export class ApiClient {
    baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    async listDevices(): Promise<Device[]> {
        const res = await fetch(`${this.baseUrl}/v1/devices`);
        if (!res.ok) throw new Error(`Failed to list devices: ${res.statusText}`);
        const data = await res.json();
        return data.devices || [];
    }

    async registerDevice(device: { device_id: string; alias: string; host: string; port: number; port_gateway: number; protocol: Protocol }): Promise<Device> {
        const res = await fetch(`${this.baseUrl}/v1/devices/${device.device_id}`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(device)
        });
        if (!res.ok) throw new Error(`Failed to register device: ${res.statusText}`);
        const data = await res.json();
        return data.device;
    }

    async updateDevice(deviceId: string, status: DeviceStatus): Promise<void> {
        const res = await fetch(`${this.baseUrl}/v1/devices/${deviceId}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ device_status: status })
        });
        if (!res.ok) throw new Error(`Failed to update device: ${res.statusText}`);
    }

    async streamDiagnostics(deviceId: string, onData: (data: DiagnosticsResponse) => void, onError: (err: any) => void): Promise<() => void> {
        const controller = new AbortController();

        try {
            const response = await fetch(`${this.baseUrl}/v1/diagnostics/${deviceId}/stream`, {
                signal: controller.signal
            });

            if (!response.ok) throw new Error(`Stream failed: ${response.statusText}`);
            if (!response.body) throw new Error('No response body');

            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';

            (async () => {
                try {
                    while (true) {
                        const { done, value } = await reader.read();
                        if (done) break;

                        buffer += decoder.decode(value, { stream: true });
                        const lines = buffer.split('\n');
                        buffer = lines.pop() || '';

                        for (const line of lines) {
                            if (line.trim()) {
                                try {
                                    const data = JSON.parse(line);
                                    onData(data.result || data);
                                } catch (e) {
                                    console.warn('Failed to parse stream line:', line, e);
                                }
                            }
                        }
                    }
                } catch (e: any) {
                    if (e.name !== 'AbortError') {
                        onError(e);
                    }
                }
            })();
        } catch (e) {
            onError(e);
        }
        return () => controller.abort();
    }
}

class AppState {
    client = $state<ApiClient | null>(null);
    devices = $state<Device[]>([]);
    isOutageActive = $state(false);
    selectedDeviceId = $state<string | null>(null);
    diagnosticsMap = $state<Record<string, Diagnostics>>({});
    diagnosticsHistory = $state<Record<string, { timestamp: number; cpu: number; memory: number }[]>>({});

    get selectedDevice() {
        return this.devices.find(d => d.device_id === this.selectedDeviceId) || null;
    }

    setClient(client: ApiClient) {
        this.client = client;
    }

    setDevices(devices: Device[]) {
        this.devices = devices;
    }

    selectDevice(deviceId: string) {
        this.selectedDeviceId = deviceId;
    }

    addDiagnosticData(deviceId: string, data: Diagnostics) {
        this.diagnosticsMap[deviceId] = data;

        if (!this.diagnosticsHistory[deviceId]) {
            this.diagnosticsHistory[deviceId] = [];
        }
        const newPoint = {
            timestamp: Date.now(),
            cpu: data.cpu_usage,
            memory: data.memory_usage
        };
        const history = this.diagnosticsHistory[deviceId];
        history.push(newPoint);
        if (history.length > 50) {
            history.shift();
        }
    }
}

export const appState = new AppState();
