use reqwest::Client;
use serde::{Deserialize, Serialize};
use thiserror::Error;

#[derive(Error, Debug)]
pub enum Error {
    #[error("request failed: {0}")]
    Reqwest(#[from] reqwest::Error),
}

pub struct HistoricalClient {
    client: Client,
    base: String,
}

#[derive(Debug, Deserialize)]
pub struct Block {
    pub slot: u64,
    pub blockhash: String,
    pub parent_slot: u64,
    pub block_time: i64,
    pub height: u64,
}

impl HistoricalClient {
    pub fn new(base: impl Into<String>) -> Self {
        Self {
            client: Client::new(),
            base: base.into(),
        }
    }
    pub async fn get_block(&self, slot: u64, commitment: &str) -> Result<Block, Error> {
        let url = format!("{}/block/{}?commitment={}", self.base, slot, commitment);
        let blk = self.client.get(&url).send().await?.json::<Block>().await?;
        Ok(blk)
    }
    pub async fn get_transaction(&self, signature: &str, commitment: &str) -> Result<serde_json::Value, Error> {
        let url = format!("{}/tx/{}?commitment={}", self.base, signature, commitment);
        let tx = self.client.get(&url).send().await?.json().await?;
        Ok(tx)
    }
}