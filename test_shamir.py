import unittest
import random
import pyshamir

def split_secret(secret, total_shares, threshold):
    shares = pyshamir.split(bytes(secret, 'utf-8'), int(total_shares), int(threshold))
    encoded_shares = [f"{i+1}-{share.hex()}" for i, share in enumerate(shares)]
    return ','.join(encoded_shares)

def restore_secret(shares):
    shares_list = shares.split(',')
    decoded_shares = [bytes.fromhex(share.split('-')[1]) for share in shares_list]
    secret = pyshamir.combine(decoded_shares)
    return secret.decode('utf-8')

def shuffle_shares(encoded_shares):
    shares = encoded_shares.split(',')
    random.shuffle(shares)
    return ','.join(shares)

def get_n_shares(encoded_shares, new_count):
    shares = encoded_shares.split(',')
    return ','.join(shares[:new_count])

def remove_indexes_from_shares(encoded_shares):
    shares = encoded_shares.split(',')
    return ','.join(share.split('-')[1] for share in shares)

def add_extra_shares(encoded_shares):
    extra_share = "4-" + bytes.hex(b"extra_share_data")
    return encoded_shares + "," + extra_share

class TestShamir(unittest.TestCase):
    def test_split_and_restore_secret(self):
        test_cases = [
            {
                "name": "simple word",
                "secret": "my_secret",
                "threshold": 2,
                "total_shares": 3,
            },
            {
                "name": "simple sentence",
                "secret": "god of mices ate my food and then he cried",
                "threshold": 2,
                "total_shares": 6,
            },
            {
                "name": "long sentence",
                "secret": "god of mices ate my food and then he cried while other mices danced around him",
                "threshold": 6,
                "total_shares": 6,
            },
            {
                "name": "BIP-39",
                "secret": "pen aunt text rotate donate sock shield pottery cloud toy tank sibling parrot oblige agent egg october angle short wolf survey frequent autumn desert",
                "threshold": 3,
                "total_shares": 6,
            },
        ]

        for tc in test_cases:
            with self.subTest(tc["name"]):
                secret = tc["secret"]
                threshold = tc["threshold"]
                total_shares = tc["total_shares"]

                print(tc)

                encoded_shares = split_secret(secret, total_shares, threshold)
                restored_secret = restore_secret(encoded_shares)
                self.assertEqual(secret, restored_secret)

                shuffled_shares = shuffle_shares(encoded_shares)
                self.assertNotEqual(encoded_shares, shuffled_shares)

                restored_secret = restore_secret(shuffled_shares)
                self.assertEqual(secret, restored_secret)

                for i in range(1, total_shares):
                    if i < threshold:
                        insufficient_shares = get_n_shares(encoded_shares, i)
                        with self.assertRaises(ValueError):
                            restore_secret(insufficient_shares)
                    else:
                        sufficient_shares = get_n_shares(encoded_shares, i)
                        restored_secret = restore_secret(sufficient_shares)
                        self.assertEqual(secret, restored_secret)

                incomplete_shares = encoded_shares[:len(encoded_shares)-2]
                with self.assertRaises(ValueError):
                    restore_secret(incomplete_shares)

                extra_shares = add_extra_shares(encoded_shares)
                with self.assertRaises(ValueError):
                    restore_secret(extra_shares)

                with self.assertRaises(IndexError):
                    restore_secret("1624ghjsgd762")

                with self.assertRaises(IndexError):
                    restore_secret(remove_indexes_from_shares(encoded_shares))

    def test_incorrect_split(self):
        test_cases = [
            {
                "name": "zero threshold",
                "threshold": 0,
                "total_shares": 3,
            },
            {
                "name": "too low threshold",
                "threshold": 1,
                "total_shares": 3,
            },
            {
                "name": "threshold > total shares",
                "threshold": 3,
                "total_shares": 2,
            },
            {
                "name": "both zero",
                "threshold": 0,
                "total_shares": 0,
            },
        ]

        for tc in test_cases:
            with self.subTest(tc["name"]):
                threshold = tc["threshold"]
                total_shares = tc["total_shares"]
                with self.assertRaises(ValueError):
                    split_secret("my_secret", total_shares, threshold)

if __name__ == "__main__":
    unittest.main()
