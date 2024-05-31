import sys
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

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 shamir.py <command> <args>")
        sys.exit(1)

    command = sys.argv[1]
    if command == "split":
        if len(sys.argv) != 5:
            print("Usage: python3 shamir.py split <secret> <threshold> <total_shares>")
            sys.exit(1)
        secret = sys.argv[2]
        threshold = sys.argv[3]
        total_shares = sys.argv[4]

        if threshold > total_shares:
            print("Threshold cannot be higher than total shares")
            sys.exit(1)
        print(split_secret(secret, total_shares, threshold))
    elif command == "restore":
        if len(sys.argv) != 3:
            print("Usage: python3 shamir.py restore <shares>")
            sys.exit(1)
        shares = sys.argv[2]
        print(restore_secret(shares))
    else:
        print("Invalid command. Use 'split' or 'restore'")
        sys.exit(1)

if __name__ == "__main__":
    main()