import axios, { AxiosInstance } from 'axios';

export interface Block {
  slot: number;
  blockhash: string;
  parentSlot: number;
  blockTime: number;
  height: number;
}

export interface Transaction {
  signature: string;
  slot: number;
  blockTime: number;
  signer: string;
  fee: number;
  computeUnits: number;
  err?: string;
}

export class HistoricalClient {
  private http: AxiosInstance;
  constructor(baseURL: string) {
    this.http = axios.create({ baseURL, timeout: 10000 });
  }
  async getBlock(slot: number, commitment = 'finalized'): Promise<Block> {
    const { data } = await this.http.get(`/block/${slot}?commitment=${commitment}`);
    return data;
  }
  async getTransaction(signature: string, commitment = 'finalized'): Promise<Transaction> {
    const { data } = await this.http.get(`/tx/${signature}?commitment=${commitment}`);
    return data;
  }
  async getSignaturesForAddress(address: string, opts?: { limit?: number; before?: string; until?: string }) {
    const params = new URLSearchParams();
    if (opts?.limit) params.set('limit', String(opts.limit));
    if (opts?.before) params.set('before', opts.before);
    if (opts?.until) params.set('until', opts.until);
    const { data } = await this.http.get(`/sigs/${address}?${params.toString()}`);
    return data;
  }
}
export default HistoricalClient;