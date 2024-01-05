# About FlowerPot and Petals

Version: {{version}}

Flowerpot is a user-friendly graphical interface designed to facilitate the installation and deployment of a "Petals Swarm" peer.

🌐 Explore the [Petals.dev](https://petals.dev) website for comprehensive information.

By launching a "Petals" peer on your machine, you contribute to the community's ability to run resource-intensive machine learning models that exceed the capacity of individual consumer machines. Thanks to passionate developers like you, various tools such as chatbots, story generators, translators, and more can be tested and developed.

❤️ Your invaluable contribution is sincerely appreciated from the bottom of our hearts ❤️!

In the Petals network, your computer acts as a peer, sharing a portion of its GPU (memory and computational power) with others.

Rest assured that **Petals users cannot execute any malicious code on your computer**.

🌐 Refer to [this page](https://github.com/bigscience-workshop/petals/wiki/Security,-privacy,-and-AI-safety) for answers to security-related questions.

---

## Auto Stop Feature

Flowerpot also provides an auto-start feature, allowing you to commence GPU sharing upon logging into Gnome, KDE, or any other window management system.

You have the flexibility to specify processes that, when detected, will automatically stop the server. For example, if you use GPU-intensive applications like Blender or "Grand Theft Auto," simply add "blender, grand theft auto" to the "Stop on Process" settings entry. This ensures that the peer is halted, and memory is released whenever any "blender" process is running.

---

## System Requirements

"Petals" is compatible with 🐍 Python versions up to 3.11, and Flowerpot seamlessly locates it on your local computer. The process involves:

- Installing "petals" within a virtual environment
- Initiating the server when you decide to share your GPU

---

## Acknowledgments to "Fyne" and "Go"

This user interface is developed using the [Fyne](https://fyne.io) library. A heartfelt thanks once again to the Fyne community for their valuable contributions.

---

## Authors

- Patrice Ferlet, [blog (in French)](https://www.metal3d.org) - Main Developer

