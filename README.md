# GoPerformancetests

## Übersicht / Overview
Dieses Repository enthält Benchmarks zur Analyse der Performance von verschiedenen Thread-Abstraktionen in der Programmiersprache Go. Die Benchmarks wurden im Rahmen der Masterarbeit "Leistungsanalyse von Thread-Abstraktionen in verschiedenen Programmiersprachen" entwickelt.

This repository contains benchmarks for analyzing the performance of different thread abstractions in the Go programming language. The benchmarks were developed as part of the master's thesis "Performance Analysis of Thread Abstractions in Different Programming Languages."

## Implementierte Benchmarks / Implemented Benchmarks
- **Mergesort Benchmark:** Vergleich der Performance von Parallelisierung mit Goroutinen
- **Bank Benchmark:** Simulation eines Bank-Systems mit verschiedenen Nebenläufigkeitsmodellen
  <br/><br/>
- **Mergesort Benchmark:** Comparison of parallelization performance using Goroutines
- **Bank Benchmark:** Simulation of a banking system using different concurrency models

## Nutzung / Usage
### Voraussetzungen / Prerequisites
- Go 1.23 oder neuer
- Docker zur Ausführung des Bank Benchmarks
  <br/><br/>
- Go 1.23 or later
- Docker for running the Bank Benchmark

### Installation / Installation
1. Repository klonen / Clone the repository:
   ```sh
   git clone https://github.com/ermerp/GoPerformancetests.git
   cd GoPerformancetests
   ```
2. Projekt bauen / Build the project:
   - **Windows (Git Bash):**
     ```sh
     ./build.sh
     docker build -t bank-go .
     ```
   - **Unix:**
     ```sh
     chmod +x build.sh
     ./build.sh
     docker build -t bank-go .
     ```

### Benchmarks ausführen / Run Benchmarks
- **Mergesort Benchmark:**
   - **Windows:**
     ```sh
     .\mergesort\mergesortGo.exe
     ```
   - **Unix:**
     ```sh
     ./mergesort/mergesortGo
     ```
- **Bank Benchmark mit Docker starten / Start the Bank Benchmark with Docker:**
   ```sh
   docker-compose up
   ```

## Monitoring & Konfiguration / Monitoring & Configuration
Für eine einfache Ausführung und Konfiguration der Benchmarks wird das Repository [Monitoring](https://github.com/ermerp/Monitoring) empfohlen.

For easy execution and configuration of the benchmarks, the repository [Monitoring](https://github.com/ermerp/Monitoring) is recommended.

## Ergebnisse / Results
Die Ergebnisse der Benchmarks werden in der Masterarbeit ausführlich analysiert. Eine detaillierte Analyse kann in der vollständigen Arbeit eingesehen werden:

The results of the benchmarks are analyzed in detail in the master's thesis. A full analysis can be found in the complete thesis:

[Download Thesis PDF](https://github.com/ermerp/masterthesis/releases/latest/download/performancevergleich.pdf)

## Lizenz / License
Dieses Projekt steht unter einer [Creative Commons Namensnennung 4.0 Lizenz](https://creativecommons.org/licenses/by/4.0/deed.de).

This project is licensed under a [Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/deed.en).
