# AspireToExpandServer

**Aspire To Expand** is an educational app used for learning English. This package is specifically for the backend services related to the app functionality, such as handling authentication and fetching resources.

_An educational platform for teachers who want a trustworthy and well built desktop application for teaching their students._ 👩‍🏫👨‍🏫💻📚

## Overview

This app provides teachers with the ability manage existing and new student data, securely accept payment for classes via Stripe and Paypal, chat via instant messaging, create and share documents via Google Drive/Docs and S3, use tools such as Google Translate and ChatGPT, conduct one-on-one classes with video/audio/screen sharing capabilities and fullscreen view (_NOTE_: support for up to 4 participants, i.e. 1 teacher and 3 students is currently being built), sending push notifications to individual students or to multiple students via student segments (i.e. send notifications to all students of a certain level, or to all students in a particular timezone, etc.) configure games for students, and more.

## Table of Contents

- [Contact info](#contact-info)
- [Features](#features)
- [Demo](#demo)
- [Installation for development](#installation-for-development)
- [Download production app](#download-production-app)
- [Usage](#usage)
- [Technologies](#technologies)
- [License](#license)

## Contact info

**Developed by:** Mitchell Wintrow

**Email**: owner@winapps.io, mitchellwintrow@gmail.com

## Features

- 🔒 Secure user authentication and authorization, as well as Google SSO
- 💳 Secure payment and PCI DSS compliance via Stripe and Paypal for classes _(Under development)_
- 🌍 Multi-language support _(Currently supports Ukrainian, Russian, and German. Support is being developed for Spanish, French, Hebrew, Arabic, Chinese, Japanese, Portugeuse, Vietnamese, and Korean. Please reach out for support with any languages not mentioned and we'll see what we can do for you!)_
- 🎥 Encrypted video conferencing and screen sharing with up to 4 partipants via WebRTC
- 💬 Instant messaging with students via WebSockets
- 🗂️ Integration with S3 (AWS) _(Under development)_
- 📂 Integration with Google Drive/Docs _(Under development)_
- ㊙️ Integration with Google Translate _(Under development)_
- 🤖 Integration with ChatGPT _(Under development)_
- 🎮 Gamification for students
- 🧍‍♂️🧍‍♀️ Student segmentation _(Under development)_
- 🔔 Push notifications _(Under development)_

## Demo

![gifDemo](https://winapps-solutions-llc.s3.us-west-2.amazonaws.com/products/aspire-with-alina/teacherDesktopDemo.gif)

**Old Demo Below: Added to showcase Video calling**

![oldGifDemo](./demoTeacherDesktop.mov.gif)

More demos will be added over time...

## Installation for development

**For Mac/Linux**

```bash
# Step 1: Clone the repo
git clone https://github.com/mrrobotisreal/AspireToExpandServer.git

# Step 2: Enter project directory and build the executable
cd AspireToExpandServer
go build main.go

# Step 3: Create required environment variables
export AWS_ACCESS_KEY_ID=<your-AWS_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<your-AWS_SECRET_ACCESS_KEY>
export AWS_REGION=<your-AWS_REGION>

# Step 4: Run the executable (or go server)
./main
# or
go run main.go
```

**For Windows (Powershell)**

```powershell
# Step 1: Clone the repo
git clone https://github.com/mrrobotisreal/AspireToExpandServer.git

# Step 2: Enter project directory and build the executable
cd AspireToExpandServer
go build main.go

# Step 3: Create required environment variables
$Env:AWS_ACCESS_KEY_ID = <your-AWS_ACCESS_KEY_ID>
$Env:AWS_SECRET_ACCESS_KEY = <your-AWS_SECRET_ACCESS_KEY>
$Env:AWS_REGION = <your-AWS_REGION>

# Step 4: Run the executable (or go server)
./main
# or
go run main.go
```

**For Windows (Command Line)**

```cmd
rem Step 1: Clone the repo
git clone https://github.com/mrrobotisreal/AspireToExpandServer.git

rem Step 2: Enter project directory and build the executable
cd AspireToExpandServer
go build main.go

rem Step 3: Create required environment variables
set AWS_ACCESS_KEY_ID=<your-AWS_ACCESS_KEY_ID>
set AWS_SECRET_ACCESS_KEY=<your-AWS_SECRET_ACCESS_KEY>
set AWS_REGION=<your-AWS_REGION>

# Step 4: Run the executable (or go server)
./main
# or
go run main.go
```

## Download production app

A download link to the finalized production app version will be posted here...

## Usage

**NOTE**: In order to first log in to this app, you must manually insert a teacher account into the MongoDB database on your instance of the AspireToExpandServer with the following properties:

```mongosh
db.students.insertOne({
  "teacherID":"<your-desired-teacher-id>",
  "emailaddress":"<your-email-address>",
  "password":"<your-password>"
})
```

## Technologies

![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?logo=mongodb&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-%23FF9900.svg?logo=amazon-web-services&logoColor=white)
![Git](https://img.shields.io/badge/Git-F05032?logo=git&logoColor=fff)
![macOS](https://img.shields.io/badge/macOS-000000?logo=apple&logoColor=F0F0F0)
![Pop!_OS](https://img.shields.io/badge/Pop!__OS-48B9C7?logo=popos&logoColor=fff)
![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?logo=ubuntu&logoColor=white)
![Windows](https://custom-icon-badges.demolab.com/badge/Windows-0078D6?logo=windows11&logoColor=white)
![Neovim](https://img.shields.io/badge/Neovim-57A143?logo=neovim&logoColor=fff)
![Visual Studio Code](https://custom-icon-badges.demolab.com/badge/Visual%20Studio%20Code-0078d7.svg?logo=vsc&logoColor=white)

[![WinApps Solutions LLC](https://img.shields.io/badge/WinApps-%232f56a0.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABwAAAAcCAMAAABF0y+mAAAACXBIWXMAAAk6AAAJOgHwZJJKAAACB1BMVEVHcEz////+/v79/f3+/v7+/v7////////////7+/v+/v78/Pz+/v79/f39/f39/f37+/v8/Pz9/f3+/v7+/v7+/v7+/v79/f38/Pz+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v739/f9/f34+Pj5+fn5+fn5+fn7+/v7+/v7+/v7+/v7+/v8/Pz8/Pzl5eX8/Pz9/f39/f3x8fH8/Pz+/v7+/v7+/v79/f39/f3+/v7+/v79/f39/f39/f39/f39/f39/f38/Pz+/v7+/v719fX19fX+/v7+/v78/Pz9/f39/f38/Pz9/f3p6en39/f29vb39/f39/f39/f39/f39/f+/v66urr+/v79/f38/Pz9/f39/f38/Pz9/f35+fn5+fn5+fn4+Pj4+Pj+/v77+/v7+/vu7u77+/v6+vr7+/v7+/v7+/v6+vrLy8v6+vr6+vru7u7j4+P7+/v8/Pz8/Pz9/f39/f39/f38/Pz8/Pzx8fH+/v7+/v7a2trz8/Pb29v+/v7+/v7+/v7n5+fz8/P+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v79/f39/f39/f3+/v7+/v7+/v7+/v7+/v79/f39/f39/f39/f39/f38/Pz9/f39/f39/f39/f39/f39/f39/f3+/v7+/v7+/v7+/v7+/v48I/X5AAAArHRSTlMAA/oDBfn99wEBvZfRSHtYAnBqbuvI48KJ2ebsxfLdvqJ/E2kBHSAPPTE7DxNeNwVsT1sGJkD06byZwc6ioce5kIyBurAMC0ddE39lZ0QDGx0XGBURBiMBBiNGUgYZBCQJGRIXdnAKATcpRGIeFgEyLwQDODUsd3gvMCAIUzEDDAF6fV8BCoqc5a+rzJCPtpue1vGGnbLV27WJ4J+3z8zSiu2wpbq0q6mh0t64yjchRwAABL1pVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0n77u/JyBpZD0nVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkJz8+Cjx4OnhtcG1ldGEgeG1sbnM6eD0nYWRvYmU6bnM6bWV0YS8nPgo8cmRmOlJERiB4bWxuczpyZGY9J2h0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMnPgoKIDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PScnCiAgeG1sbnM6QXR0cmliPSdodHRwOi8vbnMuYXR0cmlidXRpb24uY29tL2Fkcy8xLjAvJz4KICA8QXR0cmliOkFkcz4KICAgPHJkZjpTZXE+CiAgICA8cmRmOmxpIHJkZjpwYXJzZVR5cGU9J1Jlc291cmNlJz4KICAgICA8QXR0cmliOkNyZWF0ZWQ+MjAyNS0wMi0xMDwvQXR0cmliOkNyZWF0ZWQ+CiAgICAgPEF0dHJpYjpFeHRJZD5kY2RmNDhiOC03ODRiLTQ3OWEtYjA1ZS05MjNkM2I1MzM1ZmI8L0F0dHJpYjpFeHRJZD4KICAgICA8QXR0cmliOkZiSWQ+NTI1MjY1OTE0MTc5NTgwPC9BdHRyaWI6RmJJZD4KICAgICA8QXR0cmliOlRvdWNoVHlwZT4yPC9BdHRyaWI6VG91Y2hUeXBlPgogICAgPC9yZGY6bGk+CiAgIDwvcmRmOlNlcT4KICA8L0F0dHJpYjpBZHM+CiA8L3JkZjpEZXNjcmlwdGlvbj4KCiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0nJwogIHhtbG5zOmRjPSdodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyc+CiAgPGRjOnRpdGxlPgogICA8cmRmOkFsdD4KICAgIDxyZGY6bGkgeG1sOmxhbmc9J3gtZGVmYXVsdCc+Q29weSBvZiBXaW5BcHBzX0JhZGdlX1doaXRlIC0gMTwvcmRmOmxpPgogICA8L3JkZjpBbHQ+CiAgPC9kYzp0aXRsZT4KIDwvcmRmOkRlc2NyaXB0aW9uPgoKIDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PScnCiAgeG1sbnM6cGRmPSdodHRwOi8vbnMuYWRvYmUuY29tL3BkZi8xLjMvJz4KICA8cGRmOkF1dGhvcj5NaXRjaCBXaW50cm93PC9wZGY6QXV0aG9yPgogPC9yZGY6RGVzY3JpcHRpb24+CgogPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9JycKICB4bWxuczp4bXA9J2h0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8nPgogIDx4bXA6Q3JlYXRvclRvb2w+Q2FudmEgZG9jPURBR2V2a2dwblBBIHVzZXI9VUFGQXVyOXdYaDAgYnJhbmQ9QkFGQXVwTnhiSm8gdGVtcGxhdGU9PC94bXA6Q3JlYXRvclRvb2w+CiA8L3JkZjpEZXNjcmlwdGlvbj4KPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KPD94cGFja2V0IGVuZD0ncic/PupbsSAAAAIcSURBVCgVVVIHWxNBEH13Ry7BJHIBkpACCSJKVUBIAlgAsfeO7eyKXexiwdh7V1TsXXk/0tk08d337c1OeVN2gDwsV07SpuZVhb8B1Lk9nnn7dgBWQZsVDKwK0hEKk46emRDPCTBQz8il/ThweT75fsp/sRZqmTqWdf5FfrgyIU7Ej1GgqNjQigAP+QlawWyjpXIn7BxZK03HNBQXrCOfOzJyw5do74KrVeTZQk0WBn0Hz32NRb/pJIOIk/15Xun50HfRKui6k3Vp8neBFBWiNZVNDp1tZeSf1VmrC25lcoxfu/5DIsldt8jxvlKV1cBhOsnq6eI6a5JE+yDn7fOJNarg3eXknfa1cX/8CG4sWnwBw06TLYj9VDWVksklR1WdocZMoqU0E7Xo4LDcUqQXnXRKncuApoquhMkemUf5ZGAgxLvAPalSZxqNHmFgYJ2Mq/++pAyzBFhOU9pYgRph0GW0toYHYeBmgJvhwmwJSK5Hb4ZhDjQDD5ObAD83qolvqTQZx1bpVKf0aONRlRSUFlpbvufSbh9eCMPLldiG7ewU4x52A7aFi6NC+XjodFmrKFUTxyGdNvOJXAYwgxIbqFGj2uD1062SWfDxqVc0OPFMDVJ/FesaJVNDmfd2oaGbkddvqsua50rGLNwnVaBAZli/dywYfdeGklRElvPt2Kl/y2vkl108m84Mti8UkzxJHpatWZZmIecmksJfW1eka4cCl0kAAAAASUVORK5CYII=&style=flat&labelColor=%232f56a0)](https://winapps.io/products)
[![Mitchell Wintrow](https://img.shields.io/badge/Mitchell_Wintrow-%23ff6f00.svg?logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABwAAAAcCAMAAABF0y+mAAAACXBIWXMAAAk6AAAJOgHwZJJKAAACB1BMVEVHcEz////+/v79/f3+/v7+/v7////////////7+/v+/v78/Pz+/v79/f39/f39/f37+/v8/Pz9/f3+/v7+/v7+/v7+/v79/f38/Pz+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v739/f9/f34+Pj5+fn5+fn5+fn7+/v7+/v7+/v7+/v7+/v8/Pz8/Pzl5eX8/Pz9/f39/f3x8fH8/Pz+/v7+/v7+/v79/f39/f3+/v7+/v79/f39/f39/f39/f39/f39/f38/Pz+/v7+/v719fX19fX+/v7+/v78/Pz9/f39/f38/Pz9/f3p6en39/f29vb39/f39/f39/f39/f39/f+/v66urr+/v79/f38/Pz9/f39/f38/Pz9/f35+fn5+fn5+fn4+Pj4+Pj+/v77+/v7+/vu7u77+/v6+vr7+/v7+/v7+/v6+vrLy8v6+vr6+vru7u7j4+P7+/v8/Pz8/Pz9/f39/f39/f38/Pz8/Pzx8fH+/v7+/v7a2trz8/Pb29v+/v7+/v7+/v7n5+fz8/P+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v7+/v79/f39/f39/f3+/v7+/v7+/v7+/v7+/v79/f39/f39/f39/f39/f38/Pz9/f39/f39/f39/f39/f39/f39/f3+/v7+/v7+/v7+/v7+/v48I/X5AAAArHRSTlMAA/oDBfn99wEBvZfRSHtYAnBqbuvI48KJ2ebsxfLdvqJ/E2kBHSAPPTE7DxNeNwVsT1sGJkD06byZwc6ioce5kIyBurAMC0ddE39lZ0QDGx0XGBURBiMBBiNGUgYZBCQJGRIXdnAKATcpRGIeFgEyLwQDODUsd3gvMCAIUzEDDAF6fV8BCoqc5a+rzJCPtpue1vGGnbLV27WJ4J+3z8zSiu2wpbq0q6mh0t64yjchRwAABL1pVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0n77u/JyBpZD0nVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkJz8+Cjx4OnhtcG1ldGEgeG1sbnM6eD0nYWRvYmU6bnM6bWV0YS8nPgo8cmRmOlJERiB4bWxuczpyZGY9J2h0dHA6Ly93d3cudzMub3JnLzE5OTkvMDIvMjItcmRmLXN5bnRheC1ucyMnPgoKIDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PScnCiAgeG1sbnM6QXR0cmliPSdodHRwOi8vbnMuYXR0cmlidXRpb24uY29tL2Fkcy8xLjAvJz4KICA8QXR0cmliOkFkcz4KICAgPHJkZjpTZXE+CiAgICA8cmRmOmxpIHJkZjpwYXJzZVR5cGU9J1Jlc291cmNlJz4KICAgICA8QXR0cmliOkNyZWF0ZWQ+MjAyNS0wMi0xMDwvQXR0cmliOkNyZWF0ZWQ+CiAgICAgPEF0dHJpYjpFeHRJZD5kY2RmNDhiOC03ODRiLTQ3OWEtYjA1ZS05MjNkM2I1MzM1ZmI8L0F0dHJpYjpFeHRJZD4KICAgICA8QXR0cmliOkZiSWQ+NTI1MjY1OTE0MTc5NTgwPC9BdHRyaWI6RmJJZD4KICAgICA8QXR0cmliOlRvdWNoVHlwZT4yPC9BdHRyaWI6VG91Y2hUeXBlPgogICAgPC9yZGY6bGk+CiAgIDwvcmRmOlNlcT4KICA8L0F0dHJpYjpBZHM+CiA8L3JkZjpEZXNjcmlwdGlvbj4KCiA8cmRmOkRlc2NyaXB0aW9uIHJkZjphYm91dD0nJwogIHhtbG5zOmRjPSdodHRwOi8vcHVybC5vcmcvZGMvZWxlbWVudHMvMS4xLyc+CiAgPGRjOnRpdGxlPgogICA8cmRmOkFsdD4KICAgIDxyZGY6bGkgeG1sOmxhbmc9J3gtZGVmYXVsdCc+Q29weSBvZiBXaW5BcHBzX0JhZGdlX1doaXRlIC0gMTwvcmRmOmxpPgogICA8L3JkZjpBbHQ+CiAgPC9kYzp0aXRsZT4KIDwvcmRmOkRlc2NyaXB0aW9uPgoKIDxyZGY6RGVzY3JpcHRpb24gcmRmOmFib3V0PScnCiAgeG1sbnM6cGRmPSdodHRwOi8vbnMuYWRvYmUuY29tL3BkZi8xLjMvJz4KICA8cGRmOkF1dGhvcj5NaXRjaCBXaW50cm93PC9wZGY6QXV0aG9yPgogPC9yZGY6RGVzY3JpcHRpb24+CgogPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9JycKICB4bWxuczp4bXA9J2h0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC8nPgogIDx4bXA6Q3JlYXRvclRvb2w+Q2FudmEgZG9jPURBR2V2a2dwblBBIHVzZXI9VUFGQXVyOXdYaDAgYnJhbmQ9QkFGQXVwTnhiSm8gdGVtcGxhdGU9PC94bXA6Q3JlYXRvclRvb2w+CiA8L3JkZjpEZXNjcmlwdGlvbj4KPC9yZGY6UkRGPgo8L3g6eG1wbWV0YT4KPD94cGFja2V0IGVuZD0ncic/PupbsSAAAAIcSURBVCgVVVIHWxNBEH13Ry7BJHIBkpACCSJKVUBIAlgAsfeO7eyKXexiwdh7V1TsXXk/0tk08d337c1OeVN2gDwsV07SpuZVhb8B1Lk9nnn7dgBWQZsVDKwK0hEKk46emRDPCTBQz8il/ThweT75fsp/sRZqmTqWdf5FfrgyIU7Ej1GgqNjQigAP+QlawWyjpXIn7BxZK03HNBQXrCOfOzJyw5do74KrVeTZQk0WBn0Hz32NRb/pJIOIk/15Xun50HfRKui6k3Vp8neBFBWiNZVNDp1tZeSf1VmrC25lcoxfu/5DIsldt8jxvlKV1cBhOsnq6eI6a5JE+yDn7fOJNarg3eXknfa1cX/8CG4sWnwBw06TLYj9VDWVksklR1WdocZMoqU0E7Xo4LDcUqQXnXRKncuApoquhMkemUf5ZGAgxLvAPalSZxqNHmFgYJ2Mq/++pAyzBFhOU9pYgRph0GW0toYHYeBmgJvhwmwJSK5Hb4ZhDjQDD5ObAD83qolvqTQZx1bpVKf0aONRlRSUFlpbvufSbh9eCMPLldiG7ewU4x52A7aFi6NC+XjodFmrKFUTxyGdNvOJXAYwgxIbqFGj2uD1062SWfDxqVc0OPFMDVJ/FesaJVNDmfd2oaGbkddvqsua50rGLNwnVaBAZli/dywYfdeGklRElvPt2Kl/y2vkl108m84Mti8UkzxJHpatWZZmIecmksJfW1eka4cCl0kAAAAASUVORK5CYII=&style=flat)](https://winapps.io/about/mitchell-wintrow)

![Mitchell Wintrow Profile Pic](https://winapps-solutions-llc.s3.us-west-2.amazonaws.com/mitchProfilePic.png)

## License

This project is licensed under the **Creative Commons Attribution-NonCommercial-NoDerivatives 4.0 International License**.

- You can view and share this code for **non-commercial purposes** as long as proper credit is given.
- **Forking, modifications, or derivative works are not allowed.**

For the full license text, visit [Creative Commons License](https://creativecommons.org/licenses/by-nc-nd/4.0/legalcode).

---

This product is developed and owned by [WinApps (Mitchell Wintrow) ©2024](https://winapps.io)
