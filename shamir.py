import sys
import pyshamir

def split_secret(secret, total_shares, threshold):
    shares = pyshamir.split(bytes(secret, 'utf-8'), total_shares, threshold)
    encoded_shares = [f"{i+1}-{share.hex()}" for i, share in enumerate(shares)]
    return ','.join(encoded_shares)

def restore_secret(shares):
    shares_list = shares.split(',')
    decoded_shares = []
    for share in shares_list:
        parts = share.split('-', 1)
        if len(parts) != 2 or not parts[0] or not parts[1]:
            raise ValueError("Invalid share format")
        decoded_shares.append(bytes.fromhex(parts[1]))

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
        try:
            threshold = int(sys.argv[3])
            total_shares = int(sys.argv[4])
        except ValueError:
            print("Threshold and total shares must be integers")
            sys.exit(1)

        if threshold > total_shares:
            print("Threshold cannot be higher than total shares")
            sys.exit(1)
        try:
            print(split_secret(secret, total_shares, threshold))
        except ValueError as err:
            print(f"Error splitting secret: {err}")
            sys.exit(1)
    elif command == "restore":
        if len(sys.argv) != 3:
            print("Usage: python3 shamir.py restore <shares>")
            sys.exit(1)
        shares = sys.argv[2]
        try:
            print(restore_secret(shares))
        except (ValueError, UnicodeDecodeError) as err:
            print(f"Error restoring secret: {err}")
            sys.exit(1)
    else:
        print("Invalid command. Use 'split' or 'restore'")
        sys.exit(1)

if __name__ == "__main__":
    main()
