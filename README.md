![Project Logo](https://www.airchains.io/assets/logos/airchains-cosmwasm-rollup-full-logo.png)

# Overview

The CosmWasm Chain Sequencer is a groundbreaking, high-performance tool developed to enhance transaction and block management on Smart Virtual Machine (CosmWasm) chains. This innovative tool is distinguished by its integration of cutting-edge technologies and processes, tailored specifically for CosmWasm chains, to ensure streamlined, efficient, and reliable blockchain operations.

## Table of Contents

- [Overview](#overview)
    - [Table of Contents](#table-of-contents)
    - [Key Features](#key-features)
    - [Usage](#usage)
    - [Script Files Overview](#script-files-overview)
    - [License](#license)
    - [Acknowledgments](#acknowledgments)

## Key Features

- **Optimized Transaction Processing**: Utilizes advanced algorithms to effectively manage and process transactions, significantly increasing throughput and reducing latency in SVM chain environments.

- **Sophisticated Block Management**: The sequencer is equipped with refined block management capabilities, ensuring smooth and efficient block generation and propagation within the SVM ecosystem.

- **Data Integrity and Availability**: Integrates robust Data Availability (DA) processes, vital for maintaining the integrity and accessibility of data on the chain. This feature enhances the trustworthiness and transparency of the SVM network.

- **Custom Batching Techniques**: Features bespoke batching techniques that are specifically designed for SVM chains, optimizing the handling of transaction loads and enhancing overall network performance.

- **Seamless Integration with SVM Layer**: Designed to seamlessly integrate with the SVM layer, the sequencer maintains consistent performance and compatibility, reinforcing the strength and stability of SVM-based applications.

- **High Scalability and Flexibility**: Built to accommodate the growing demands of SVM chains, offering scalable solutions that adapt to varying transaction volumes and network conditions.

- **Reliability and Security**: Prioritizes reliability and security in its design, ensuring that the SVM Chain Sequencer operates with high resilience and robustness, safeguarding against potential threats and vulnerabilities.

## Usage

In order to tailor the Sequencer to better align with your specific requirements, please proceed to update key configuration parameters within the `common/constants.go` file. The following constants are crucial for the optimal functioning of the sequencer and can be adjusted to meet your operational needs:

- **BatchSize**: Modify this value to alter the batch size for transaction processing. This adjustment can optimize throughput and efficiency based on your workload.

- **BlockDelay**: Adjust this constant to set the delay between blocks check, aligning it with your network's block generation rate for synchronized operations.

- **ExecutionClientRPC**: Update this URL to connect the sequencer with your preferred execution client's RPC interface.

- **SettlementClientRPC**: Change this URL to integrate the sequencer with the desired settlement layer's RPC service.

- **KeyringDirectory**: Specify a new directory path for the keyring, ensuring secure and organized storage of cryptographic keys.

- **DaClientRPC**: Alter this URL to link the sequencer with your chosen Data Availability (DA) service's RPC endpoint.

Each of these parameters plays a critical role in the configuration and performance of the sequencer. It is recommended to carefully consider the implications of these changes to maintain optimal functionality and security of the system.

## Script Files Overview

This repository includes a set of four essential script files designed to facilitate various operational aspects of the node. Each script has a specific purpose, detailed below:

- **init.sh**: This script is responsible for setting up the initial environment. It creates all necessary folders and files required for the node's operation.

```
  sh scripts/init.sh
```

- **start.sh**: Designed to simply start the node without altering its current state or block height. Ideal for routine starts where no changes to the environment are needed.

```
  sh scripts/start.sh
```

- **restart.sh**: Used for refreshing the node's environment. It first deletes and then recreates all folders and files, followed by restarting the node to apply these changes.

```
  sh scripts/restart.sh
```

- **test.sh**: As the name suggests, this script is used for testing purposes. It contains commands and procedures specifically intended for testing the node's functionalities and performance.

```
  sh scripts/test.sh
```

> Note: before proceeding to run the sequencer, please ensure that the `init.sh` script has been executed to initialize the basic directory structure and configuration files.

_Important Security Notice Regarding init.sh Execution_

Please be aware that running the `init.sh` script necessitates the entry of your terminal password. This requirement stems from the inclusion of `sudo` commands within the script. These commands elevate privileges for certain operations, which are essential for the correct setup and configuration of the environment.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

Special thanks to the `gnark` library, an efficient and elegant toolkit for zk-SNARKs on Go. This library has been instrumental in our development process. For more information and to explore their work, visit their GitHub repository at [Consensys/gnark.](https://github.com/Consensys/gnark)
