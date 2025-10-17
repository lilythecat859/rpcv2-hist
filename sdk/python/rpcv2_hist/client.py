import httpx
from typing import Optional, List, Dict, Any

class HistoricalClient:
    def __init__(self, base_url: str, timeout: float = 10):
        self.client = httpx.Client(base_url=base_url, timeout=timeout)

    def get_block(self, slot: int, commitment: str = "finalized") -> Dict[str, Any]:
        r = self.client.get(f"/block/{slot}?commitment={commitment}")
        r.raise_for_status()
        return r.json()

    def get_transaction(self, signature: str, commitment: str = "finalized") -> Dict[str, Any]:
        r = self.client.get(f"/tx/{signature}?commitment={commitment}")
        r.raise_for_status()
        return r.json()

    def get_signatures_for_address(
        self,
        address: str,
        limit: int = 100,
        before: Optional[str] = None,
        until: Optional[str] = None,
    ) -> List[Dict[str, Any]]:
        params: Dict[str, Any] = {"limit": limit}
        if before:
            params["before"] = before
        if until:
            params["until"] = until
        r = self.client.get(f"/sigs/{address}", params=params)
        r.raise_for_status()
        return r.json()