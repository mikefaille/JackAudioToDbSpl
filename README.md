
# JackAudioToDbSpl

Convert PCM data from a JACK Audio channel to measure Sound Pressure Level (SPL) in decibels (dB).

## 🚀 Features

- Connects with JACK Audio channels.
- Prints non-zero samples for debugging.
- Extensible to add more features like SPL measurement in the future.

## 🛠️ Installation

### Prerequisites
- JACK Audio Connection Kit: For obtaining PCM data.

### Steps
1. Clone the repository:
```bash
git clone https://github.com/mikefaille/JackAudioToDbSpl
```
2. Navigate to the project directory:
```bash
cd JackAudioToDbSpl
```
3. Build the application:
```bash
go build
```

## 📋 Usage

Run the compiled binary:
```bash
./JackAudioToDbSpl
```

## 🤝 Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are greatly appreciated.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📝 License

Distributed under the MIT License. See `LICENSE` for more information.

## 📧 Contact

Michael Faille - michael -at- faille - io

## 🌟 Acknowledgements

- [go-jack](https://github.com/xthexder/go-jack): Go bindings for the JACK Audio Connection Kit.

