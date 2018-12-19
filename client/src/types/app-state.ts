import { grpc } from 'grpc-web-client';
import { Client } from './service-client';
import { HOST } from 'src/api/host';

export interface AppStateContent {
    token?: string;
    session?: string;
    user?: { name: string }

    readonly error?: string
    setError: (msg: string) => void;
    resetSession(): void;
    logout(): void;
}

export function getClient(state: AppStateContent): Client {
    return new Client(state, HOST);
}

export function getAuthentication(state: AppStateContent): grpc.Metadata | undefined {
    if (state.token) {
        const md = new Map<string, string>();
        md.set("authorization", state.token);
        return new grpc.Metadata(md);
    }

    return;
}