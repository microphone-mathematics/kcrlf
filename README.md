# kcrlf
[1.1]: http://i.imgur.com/tXSoThF.png
[1]: https://twitter.com/TobiunddasMoe
This an adaption of Emoe's adaptation of tomnomnom's kxss tool with a different output format. I didn't want to fork his whole Hacks-Repository so created my Own ;-)

All Credit for this Code goes to [Tomnomnom](https://github.com/tomnomnom/) and [Emoe](https://github.com/Emoe/)

## Output
Output Looks like this:
```
URL: https://www.**********.***/event_register.php?event=177 Param: event Payload: [%0D%0Aquasimoto%3A+has-crlf %0Aquasimoto%3A+has-crlf]
```

## Installation
To install this Tool please use the following Command:
```
go get github.com/microphone-mathematics/kcrlf
```

## Usage
To run this script use the following command:
```
echo "https://www.**********.***/event_register.php?event=177" | kcrlf
```

## Question
If you have an question you can create an Issue or ping me on [![alt text][1.1]][1]
