class AndroidMissingIDError(Exception):
    def __init__(self, transaction_id: str):
        self.transaction_id = transaction_id
        super().__init__(f"Android device_id missing for transaction: {transaction_id}")

class CampaignMissingIDError(Exception):
    def __init__(self, transaction_id: str):
        self.transaction_id = transaction_id
        super().__init__(f"Campaign ID missing for transaction: {transaction_id}")

class DuplicateTransactionError(Exception):
    def __init__(self, transaction_id: str):
        self.transaction_id = transaction_id
        super().__init__(f"Transaction {transaction_id} is already locked or processed.")