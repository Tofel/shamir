# Shamir's Secret Sharing in Python

This is a simple Python program that demonstrates Shamir's Secret Sharing, a cryptographic algorithm used to split a secret into multiple shares, where only a certain number of shares are needed to reconstruct the original secret. This can be useful for securely storing sensitive information.

## How It Works

Shamir's Secret Sharing works by splitting a secret into multiple parts (shares). You need a certain number of these shares (threshold) to reconstruct the original secret. If you have fewer than the threshold number of shares, you cannot reconstruct the secret.

### Features

- **Split a Secret**: Divide a secret into multiple shares.
- **Restore a Secret**: Reconstruct the original secret from the shares.

## Requirements

- **Python 3**: Version 3.12.0 or higher.

## Installation

This repository includes a `virtualenv` folder with all necessary dependencies, managed using a Python virtual environment (`venv`). Although it won't work out-of-the box, since all paths are relative to machine, on which it was created. It's included here just as a backup of dependencies.

To use it, run:
```sh
source virtualenv/bin/activate
```

It's recommended to setup a new environment on new machine:

```sh
rm -rf virtualenv
python3 -m venv virtualenv
source virtualenv/bin/activate  # On Windows use `virtualenv\Scripts\activate`
pip3 install -r requirements.txt
```

## Usage

### Operational Security Note

Passing secrets or shares as command-line arguments is convenient, but those values may be exposed through:

- shell history
- local process listings
- terminal scrollback or logs
- Docker runtime metadata while a container is running

If you are working with sensitive material, prefer running the local Python script on an airgapped machine and temporarily disable Bash history for the session:

```sh
unset HISTFILE
set +o history
```

After you are done, you can re-enable history:

```sh
set -o history
```

If your recovery workflow is wallet-focused, one pragmatic way to reduce risk is to treat the recovered wallet as short-lived:

1. Recover the seed on an airgapped machine.
2. Verify that it matches the expected wallet.
3. Move funds to a fresh wallet as soon as practical.
4. Stop relying on the recovered seed after the transfer is complete.

This does not prevent local exposure during the recovery session, but it can reduce the value of any leaked secret by limiting how long the restored wallet remains in active use.

This strategy also does not protect against restoring the wrong secret. If shares are corrupted or insufficient, the software may still reconstruct incorrect output. For wallet recovery, you should verify that the restored wallet matches your expected accounts before moving funds.

### Splitting a Secret

To split a secret, use the `split` command followed by the secret, the threshold number of shares needed to reconstruct the secret, and the total number of shares you want to create.

```sh
python3 shamir.py split <secret> <threshold> <total_shares>
```

**Example:**

```sh
python3 shamir.py split "mysecret" 3 5
```

Safer Bash session example:

```sh
unset HISTFILE
set +o history
python3 shamir.py split "mysecret" 3 5
set -o history
```

This will split the secret "mysecret" into 5 shares, where any 3 shares are needed to restore the secret.

### Restoring a Secret

To restore a secret, use the `restore` command followed by the encoded shares (a comma-separated string).

```sh
python3 shamir.py restore <encoded_shares>
```

**Example:**

```sh
python3 shamir.py restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
```

Safer Bash session example:

```sh
unset HISTFILE
set +o history
python3 shamir.py restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
set -o history
```

This will reconstruct the original secret from the provided shares.

## Example Workflow

1. **Split the Secret:**

    ```sh
    python3 shamir.py split "mysecret" 3 5
    ```

    Output:

    ```sh
    1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a,4-3b8bcd8912,5-7bcce3b6a2
    ```

2. **Restore the Secret:**

    Use at least 3 of the encoded shares from the split step.

    ```sh
    python3 shamir.py restore "1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a"
    ```

    Output:

    ```sh
    mysecret
    ```

## Using the Docker Image

You can also use the provided Docker image to run the Shamir's Secret Sharing program without needing to set up a Python environment on your local machine. The build command will create an image compatible with the architecture of the machine it is run on. Here’s how you can do it:

1. **Building the Docker Image:**

    ```sh
    docker build -t shamir-python:latest .
    ```

2. **Running the Docker Container:**

    To split a secret:

    ```sh
    docker run --rm --network none shamir-python:latest split <secret> <threshold> <total_shares>
    ```

    Example:

    ```sh
    unset HISTFILE
    set +o history
    docker run --rm --network none shamir-python:latest split "mysecret" 3 5
    set -o history
    ```

    To restore a secret:

    ```sh
    docker run --rm --network none shamir-python:latest restore <encoded_shares>
    ```

    Example:

    ```sh
    unset HISTFILE
    set +o history
    docker run --rm --network none shamir-python:latest restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
    set -o history
    ```

Docker note:

- `--rm` removes the container after it exits, reducing persistence of runtime metadata.
- `--network none` prevents accidental network access during sensitive operations.
- Even with these flags, secrets passed as arguments may still be visible locally while the container is running.
- For the lowest exposure, prefer the local Python script over `docker run` when handling real secrets.
- If you recover a wallet in Docker anyway, treat that wallet as temporary and sweep funds to a fresh wallet after verification.

### Loading Docker Images from Disk

If you have a Docker image saved on disk, you can load it into Docker using the following command:

```sh
docker load -i path/to/your/image.tar
```

Replace `path/to/your/image.tar` with the actual path to your Docker image file.

## Notes

- The `split` command requires a secret string, a threshold, and the total number of shares.
- The `restore` command requires the encoded shares in a specific format.
- Make sure the threshold is less than or equal to the total number of shares.
- Ensure you have Python 3.12.0 or higher installed and have activated the virtual environment using the provided `venv` setup instructions.
- If you pass secrets or shares as command-line arguments, disable Bash history first if you do not want them recorded locally.
- For wallet recovery, verify that the restored wallet is the expected one before sending funds, then move funds to a fresh wallet if you want to minimize the impact of local exposure during recovery.

## Error Handling

The program includes basic error handling for invalid inputs and incorrect usage. Make sure to follow the correct command format to avoid errors.

## License

This project is licensed under the MIT License.

# Shamir's Secret Sharing w Pythonie

To prosty program w Pythonie, który demonstruje działanie algorytmu Shamir's Secret Sharing, służącego do dzielenia sekretu na wiele części, z których tylko określona liczba jest potrzebna do odtworzenia oryginalnego sekretu. Może to być przydatne do bezpiecznego przechowywania wrażliwych informacji.

## Jak to działa

Algorytm Shamir's Secret Sharing dzieli sekret na wiele części (udziałów). Aby odtworzyć oryginalny sekret, potrzebujesz określonej liczby tych udziałów (próg). Jeśli masz mniej udziałów niż wynosi próg, nie możesz odtworzyć sekretu.

### Funkcje

- **Podział sekretu**: Podziel sekret na wiele udziałów.
- **Odtworzenie sekretu**: Odtwórz oryginalny sekret z udziałów.

## Wymagania

- **Python 3**: Wersja 3.12.0 lub wyższa.

## Instalacja

To repozytorium zawiera folder `virtualenv` ze wszystkimi niezbędnymi zależnościami, zarządzanymi przy użyciu wirtualnego środowiska Pythona (`venv`). Chociaż nie zadziała od razu, ponieważ wszystkie ścieżki są względne względem maszyny, na której został utworzony, jest dołączony jako kopia zapasowa zależności.

Aby go użyć, uruchom:
```sh
source virtualenv/bin/activate
```

Zaleca się utworzenie nowego środowiska na nowej maszynie:

```sh
rm -rf virtualenv
python3 -m venv virtualenv
source virtualenv/bin/activate  # W systemie Windows użyj `virtualenv\Scripts\activate`
pip3 install -r requirements.txt
```

## Użycie

### Podział sekretu

Aby podzielić sekret, użyj polecenia `split` z następującymi argumentami: sekret, liczba udziałów potrzebnych do odtworzenia sekretu (próg) oraz całkowita liczba udziałów.

```sh
python3 shamir.py split <secret> <threshold> <total_shares>
```

**Przykład:**

```sh
python3 shamir.py split "mojsektet" 3 5
```

To polecenie podzieli sekret "mojsektet" na 5 udziałów, z których dowolne 3 są potrzebne do odtworzenia sekretu.

### Odtworzenie sekretu

Aby odtworzyć sekret, użyj polecenia `restore` z zakodowanymi udziałami (ciąg oddzielony przecinkami).

```sh
python3 shamir.py restore <encoded_shares>
```

**Przykład:**

```sh
python3 shamir.py restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
```

To polecenie odtworzy oryginalny sekret z dostarczonych udziałów.

## Przykładowy Przebieg

1. **Podział sekretu:**

    ```sh
    python3 shamir.py split "mojsektet" 3 5
    ```

    Wynik:

    ```sh
    1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a,4-3b8bcd8912,5-7bcce3b6a2
    ```

2. **Odtworzenie sekretu:**

    Użyj co najmniej 3 z zakodowanych udziałów z kroku podziału.

    ```sh
    python3 shamir.py restore "1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a"
    ```

    Wynik:

    ```sh
    mojsektet
    ```

## Użycie obrazu Docker

Możesz również użyć dostarczonego obrazu Docker, aby uruchomić program Shamir's Secret Sharing bez konieczności konfiguracji środowiska Pythona na swoim lokalnym komputerze. Polecenie build stworzy obraz zgodny z architekturą maszyny, na której jest uruchomione. Oto jak to zrobić:

1. **Budowanie obrazu Docker:**

    ```sh
    docker build -t shamir-python:latest .
    ```

2. **Uruchamianie kontenera Docker:**

    Aby podzielić sekret:

    ```sh
    docker run --rm shamir-python:latest split <secret> <threshold> <total_shares>
    ```

    Przykład:

    ```sh
    docker run --rm shamir-python:latest split "mojsektet" 3 5
    ```

    Aby odtworzyć sekret:

    ```sh
    docker run --rm shamir-python:latest restore <encoded_shares>
    ```

    Przykład:

    ```sh
    docker run --rm shamir-python:latest restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
    ```

### Ładowanie obrazów Docker z dysku

Jeśli masz obraz Docker zapisany na dysku, możesz go załadować do Dockera, używając następującego polecenia:

```sh
docker load -i path/to/your/image.tar
```

Zastąp `path/to/your/image.tar` rzeczywistą ścieżką do pliku obrazu Docker.

## Uwagi

- Polecenie `split` wymaga podania ciągu znaków (sekretu), progu oraz całkowitej liczby udziałów.
- Polecenie `restore` wymaga podania zakodowanych udziałów w określonym formacie.
- Upewnij się, że próg jest mniejszy lub równy całkowitej liczbie udziałów.
- Upewnij się, że masz zainstalowany Python w wersji 3.12.0 lub wyższej i aktywowałeś wirtualne środowisko za pomocą podanych instrukcji `venv`.

## Obsługa błędów

Program zawiera podstawową obsługę błędów dla nieprawidłowych danych wejściowych i niewłaściwego użycia. Upewnij się, że stosujesz się do poprawnego formatu poleceń, aby uniknąć błędów.

## Licencja

Ten projekt jest licencjonowany na warunkach licencji MIT.
