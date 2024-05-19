import argparse
import random
import functools

# Define a large prime number
PRIME = 2**1280 - 1

def mod_inverse(x, prime):
    return pow(x, prime - 2, prime)

def eval_at(poly, x, prime):
    accum = 0
    for coefficient in reversed(poly):
        accum *= x
        accum += coefficient
        accum %= prime
    return accum

def random_polynomial(degree, intercept, prime):
    return [intercept] + [random.randrange(0, prime) for _ in range(degree)]

def make_shares(secret, num_shares, threshold, prime=PRIME):
    if threshold > num_shares:
        raise ValueError("Threshold cannot be greater than number of shares")
    poly = random_polynomial(threshold - 1, secret, prime)
    shares = [(i, eval_at(poly, i, prime)) for i in range(1, num_shares + 1)]
    return shares

def lagrange_interpolate(x, x_s, y_s, prime):
    def PI(vals): return functools.reduce(lambda x, y: x * y, vals, 1)
    k = len(x_s)
    assert k == len(set(x_s)), "Points must be distinct"
    nums = []
    dens = []
    for i in range(k):
        others = list(x_s)
        cur = others.pop(i)
        nums.append(PI(x - o for o in others))
        dens.append(PI(cur - o for o in others))
    den = PI(dens)
    num = sum([nums[i] * y_s[i] * den * mod_inverse(dens[i], prime)
               for i in range(k)])
    return (num % prime) * mod_inverse(den, prime) % prime

def recover_secret(shares, prime=PRIME):
    if len(shares) < 2:
        raise ValueError("At least two shares needed to recover the secret")
    x_s, y_s = zip(*shares)
    return lagrange_interpolate(0, x_s, y_s, prime)

def split_secret(secret, num_shares, threshold, prime=PRIME):
    secret_bytes = secret.encode('utf-8')
    secret_int = int.from_bytes(secret_bytes, 'big')
    shares = make_shares(secret_int, num_shares, threshold, prime)
    max_length = max(len(format(y, 'x')) for _, y in shares)
    return [(x, format(y, '0' + str(max_length) + 'x')) for x, y in shares]


def recover_secret_from_shares(shares, prime=PRIME):
    shares_int = [(x, int(y, 16)) for x, y in shares]
    secret_int = recover_secret(shares_int, prime)
    secret_bytes = secret_int.to_bytes((secret_int.bit_length() + 7) // 8, 'big')
    return secret_bytes


def main():
    parser = argparse.ArgumentParser(description="Shamir's Secret Sharing CLI")
    parser.add_argument("action", choices=["split", "restore"], help="Action to perform: 'split' or 'restore'")
    parser.add_argument("input", help="Input string for splitting or comma-separated shares for restoring")
    parser.add_argument("--num_shares", type=int, help="Number of shares to split into (required for split action)")
    parser.add_argument("--threshold", type=int, help="Minimum number of shares required to restore (required for split action)")

    args = parser.parse_args()

    if args.action == "split":
        if not args.num_shares or not args.threshold:
            parser.error("The 'split' action requires --num_shares and --threshold arguments")
        input_string = args.input.strip()
        shares = split_secret(input_string, args.num_shares, args.threshold)
        print("Shares:")
        for x, y in shares:
            print(f"{x}-{y}")

    elif args.action == "restore":
        shares_str = args.input.split(',')
        shares = [(int(share.split('-')[0]), share.split('-')[1]) for share in shares_str]
        restored_secret = recover_secret_from_shares(shares)
        try:
            print("Restored String:")
            print(restored_secret.decode('utf-8'))  # Attempt to decode as UTF-8
        except UnicodeDecodeError:
            print("Restored data is not valid UTF-8. Here's the raw data:")
            print(restored_secret)

if __name__ == "__main__":
    main()