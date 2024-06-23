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

This repository includes a `vendor` folder with all necessary dependencies, managed using a Python virtual environment (`venv`). 

To set up the environment, run:

```sh
python3 -m venv vendor
source vendor/bin/activate  # On Windows use `vendor\Scripts\activate`
pip3 install -r requirements.txt
```

## Usage

### Splitting a Secret

To split a secret, use the `split` command followed by the secret, the threshold number of shares needed to reconstruct the secret, and the total number of shares you want to create.

```sh
python3 shamir.py split <secret> <threshold> <total_shares>
```

**Example:**

```sh
python3 shamir.py split "mysecret" 3 5
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
    docker run --rm shamir-python:latest split <secret> <threshold> <total_shares>
    ```

    Example:

    ```sh
    docker run --rm shamir-python:latest split "mysecret" 3 5
    ```

    To restore a secret:

    ```sh
    docker run --rm shamir-python:latest restore <encoded_shares>
    ```

    Example:

    ```sh
    docker run --rm shamir-python:latest restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
    ```

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

## Error Handling

The program includes basic error handling for invalid inputs and incorrect usage. Make sure to follow the correct command format to avoid errors.

## License

This project is licensed under the MIT License.

# Shamir's Secret Sharing w Pythonie

To jest prosty program w języku Python, który demonstruje metodę Shamir's Secret Sharing, algorytm kryptograficzny używany do podziału sekretu na wiele części, z których tylko pewna liczba jest potrzebna do odtworzenia oryginalnego sekretu. Może to być użyteczne do bezpiecznego przechowywania wrażliwych informacji.

## Jak to działa

Shamir's Secret Sharing działa poprzez podział sekretu na wiele części (udziałów). Potrzebujesz pewnej liczby tych udziałów (próg), aby odtworzyć oryginalny sekret. Jeśli masz mniej niż wymaganą liczbę udziałów, nie możesz odtworzyć sekretu.

### Funkcje

- **Podział sekretu**: Podziel sekret na wiele udziałów.
- **Odtworzenie sekretu**: Odtwórz oryginalny sekret z udziałów.

## Wymagania

- **Python 3**: Wersja 3.12.0 lub wyższa.

## Instalacja

To repozytorium zawiera folder `vendor` ze wszystkimi niezbędnymi zależnościami, zarządzanymi przy użyciu wirtualnego środowiska Pythona (`venv`). 

Aby skonfigurować środowisko, uruchom:

```sh
python3 -m venv vendor
source vendor/bin/activate  # Na Windows użyj `vendor\Scripts\activate`
pip install -r requirements.txt
```

## Użycie

### Podział sekretu

Aby podzielić sekret, użyj komendy `split` wraz z sekretem, liczbą udziałów potrzebnych do odtworzenia sekretu oraz całkowitą liczbą udziałów, które chcesz stworzyć.

```sh
python3 shamir.py split <sekret> <próg> <całkowita_liczba_udziałów>
```

**Przykład:**

```sh
python3 shamir.py split "mójsekret" 3 5
```

To podzieli sekret "mójsekret" na 5 udziałów, z czego dowolne 3 są potrzebne do odtworzenia sekretu.

### Odtworzenie sekretu

Aby odtworzyć sekret, użyj komendy `restore` wraz z zakodowanymi udziałami (ciąg rozdzielony przecinkami).

```sh
python3 shamir.py restore <zakodowane_udzialy>
```

**Przykład:**

```sh
python3 shamir.py restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
```

To odtworzy oryginalny sekret z dostarczonych udziałów.

## Przykładowy przebieg

1. **Podział sekretu:**

    ```sh
    python3 shamir.py split "mójsekret" 3 5
    ```

    Wynik:

    ```sh
    1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a,4-3b8bcd8912,5-7bcce3b6a2
    ```

2. **Odtworzenie sekretu:**

    Użyj co najmniej 3 zakodowanych udziałów z kroku podziału.

    ```sh
    python3 shamir.py restore "1-4b6a27bcce,2-3e48c6b830,3-3f3f4c4e1a"
    ```

    Wynik:

    ```sh
    mójsekret
    ```

## Używanie obrazu Docker

Możesz również użyć dostarczonego obrazu Docker, aby uruchomić program Shamir's Secret Sharing bez konieczności konfigurowania środowiska Python na swoim lokalnym komputerze. Polecenie build utworzy obraz zgodny z architekturą maszyny, na której jest uruchamiane. Oto jak to zrobić:

1. **Budowanie obrazu Docker:**

    ```sh
    docker build -t shamir-python:latest .
    ```

2. **Uruchamianie kontenera Docker:**

    Aby podzielić sekret:

    ```sh
    docker run --rm shamir-python:latest split <sekret> <próg> <całkowita_liczba_udziałów>
    ```

    Przykład:

    ```sh
    docker run --rm shamir-python:latest split "mójsekret" 3 5
    ```

    Aby odtworzyć sekret:

    ```sh
    docker run --rm shamir-python:latest restore <zakodowane_udzialy>
    ```

    Przykład:

    ```sh
    docker run --rm shamir-python:latest restore "1-3b6a27bcce3b6d4a3e48c6b8303f3f4c4e1a3b8bcd89,2-4e1a3b8bcd89123b6a27bcce3b6d4a3e48c6b8303f3f4c4e"
    ```

### Ładowanie obrazów Docker z dysku

Jeśli masz obraz Docker zapisany na dysku, możesz załadować go do Docker przy użyciu następującego polecenia:

```sh
docker load -i sciezka/do/twojego/obrazu.tar
```

Zastąp `sciezka/do/twojego/obrazu.tar` faktyczną ścieżką do pliku obrazu Docker.

## Uwagi

- Komenda `split` wymaga podania sekretną frazę, progu oraz całkowitej liczby udziałów.
- Komenda `restore` wymaga zakodowanych udziałów w określonym formacie.
- Upewnij się, że próg jest mniejszy lub równy całkowitej liczbie udziałów.
- Upewnij się, że masz zainstalowaną wersję Python 3.12.0 lub wyższą i że aktywowałeś wirtualne środowisko przy użyciu dostarczonych instrukcji `venv`.

## Obsługa błędów

Program zawiera podstawową obsługę błędów dla nieprawidłowych danych wejściowych i niepoprawnego użycia. Upewnij się, że przestrzegasz poprawnego formatu komend, aby uniknąć błędów.

## Licencja

Ten projekt jest licencjonowany na podstawie licencji MIT.