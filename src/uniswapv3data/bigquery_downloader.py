from datetime import datetime
from typing import Optional, List

import pandas as pd
import numpy as np
from google.cloud import bigquery
import google.auth


_SWAP_TOPIC = "0xc42079f94a6350d7e6235f29174924f928cc2ac818eb64fed8004e115fbcca67"


def _uint256_to_int256(num: int) -> int:
    if (num & (1 << 255)) != 0:
        return num - (1 << 256)
    return num


def _parse_log_ints(data: str) -> List[int]:
    data = data.lower()
    assert data[:2] == "0x"

    return [int(data[2 + i : 2 + i + 64], 16) for i in range(0, len(data[2:]), 64)]


def _parse_swap_data(data: str) -> List[np.float64]:
    data_ints = _parse_log_ints(data)
    return list(
        map(
            np.float64,
            [
                _uint256_to_int256(data_ints[0]),  # amount0
                _uint256_to_int256(data_ints[1]),  # amount1
                np.float64((data_ints[2] ** 2) / 2**192),  # price
                data_ints[3],  # liquidity
                data_ints[4],  # tick
            ],
        )
    )


class BigQueryUniswapV3Downloader:
    def __init__(self, credentials_path: Optional[str] = None) -> None:
        if credentials_path is None:
            self.client = bigquery.Client()
        else:
            self.client = bigquery.Client(
                credentials=google.auth.load_credentials_from_file(credentials_path)[0]
            )

    def download_pool_data(
        self, pool_address: str, from_date: datetime, to_date: datetime
    ) -> pd.DataFrame:
        result = self.client.query(
            f"""
            SELECT
                block_number,
                transaction_index,
                log_index,
                block_timestamp,
                data
            FROM
                `bigquery-public-data.crypto_ethereum.logs`
            WHERE
                address = "{pool_address.lower()}"
                AND ARRAY_LENGTH(topics) = 3
                AND topics[OFFSET(0)] = "{_SWAP_TOPIC.lower()}"
                AND block_timestamp >= TIMESTAMP("{from_date.isoformat()}")
                AND block_timestamp <= TIMESTAMP("{to_date.isoformat()}")
            ORDER BY
                block_number ASC,
                transaction_index ASC,
                log_index ASC
        """
        ).to_dataframe()

        result[["amount0", "amount1", "price", "liquidity", "tick"]] = (
            result["data"].map(_parse_swap_data).apply(pd.Series)
        )
        return result
