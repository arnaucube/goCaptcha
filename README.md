# goCaptcha
captcha server, with own datasets, to train own machine learning AI

## 0 - Why?
Captcha systems are useful to avoid bots posting data on databases. But modern captcha systems are from enterprises to train machine learning algorithms, and monetize results.
When user answers a captcha, is training the AI from the enterprise.

This project, aims to be a self hosted captcha system, that trains own AI.
To avoid feeding AI from companies.

## 1 - How to use?
### 1.1 - Frontend
Insert this lines in the html file:
```html
    <link rel="stylesheet" href="goCaptcha.css">

    <div id="goCaptcha"></div>

    <script>//define the url where is placed the server of goCaptcha
        var goCaptchaURL = "http://127.0.0.1:3025";
    </script>
    <script src="goCaptcha.js"></script>
```

It will place the goCaptcha box in the div:

![goCaptcha](https://raw.githubusercontent.com/arnaucode/goCaptcha/master/demo01.png "goCaptcha")

### 1.2 - Backend
- Put dataset images in the folder 'imgs'.
- Configure serverConfig.json:
```json
{
    "serverIP": "127.0.0.1",
    "serverPort": "3025",
    "imgsFolder": "imgs",
    "numImgsCaptcha": 9,
    "suspiciousIPCountLimit": 2,
    "timeBan": 30
}

```

- Run MongoDB.
- Go to the folder /goCaptcha, and run:
```
./goCaptcha
```
It will show:
```
user@laptop:~/goCaptcha$ ./goCaptcha
    goCaptcha started
    dataset read
    num of dataset categories: 4
    server running
    port: 3025
```


## 2 - How to make the petitions?

### 2.1 - Get the captcha:
```
GET/ 127.0.0.1:3025/captcha
```
Server response:
```json
{
    "id": "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
    "imgs": [
        "7cf6f630-e78f-469c-85dd-2d677996fea1.png",
        "d4014318-f875-4b42-b704-4f5bf5e5e00c.png",
        "2dd69b44-903d-4e78-bb7b-f8b07877c9e5.png",
        "2954fc38-819d-40c9-ae3e-7b6fbb68ddbe.png",
        "b060f58a-d44b-4e05-b466-92aa801a2aa1.png",
        "1b838c46-b784-471e-b143-48be058c39a7.png"
    ],
    "question": "leopard",
    "date": "1502274893"
}
```

### 2.2 - User selects the images that fit in the 'question' parameter
(in this case, 'leopard')

The selection is stored in an array:
```js
    selection=[0,0,1,0,1,0];
```
Where the '1' are the images selected, in the images array order.

### 2.3 - Post the answer. The answer contains the CaptchaId, and an array with the selected images
```
POST/ 127.0.0.1:3025/answer
```
Post example:
```json
{
	"captchaid": "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
	"selection": [0,0,1,0,1,0]
}
```
Server response:
```
true
```

## 3 - How this works?

### 3.1 - Server reads dataset
First, server reads all dataset. Dataset is a directory with subdirectories, where each subdirectory contains images of one element.

For example:
```
imgs/
    leopard/
        img01.png
        img02.png
        img03.png
        ...
    laptop/
        img01.png
        img02.png
        ...
    house/
        img01.png
        img02.png
        ...
```
Then, stores all the filenames corresponding to each subdirectory. So, we have each image and to which element category is (the name of subdirectory).


### 3.2 - Server generates captcha
When server recieves a GET /captcha, generates a captcha, getting random images from the dataset.

For each captcha generated, generates two mongodb models:
```json
Captcha Model
{
    "id" : "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
    "imgs" : [
        "7cf6f630-e78f-469c-85dd-2d677996fea1.png",
        "d4014318-f875-4b42-b704-4f5bf5e5e00c.png",
        "2dd69b44-903d-4e78-bb7b-f8b07877c9e5.png",
        "2954fc38-819d-40c9-ae3e-7b6fbb68ddbe.png",
        "b060f58a-d44b-4e05-b466-92aa801a2aa1.png",
        "1b838c46-b784-471e-b143-48be058c39a7.png"
    ],
    "question" : "leopard"
}
```

```json
CaptchaSolution Model
{
    "id" : "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
    "imgs" : [
        "image_0022.jpg",
        "image_0006.jpg",
        "image_0050.jpg",
        "image_0028.jpg",
        "image_0119.jpg",
        "image_0092.jpg"
    ],
    "imgssolution" : [
        "camera",
        "camera",
        "laptop",
        "crocodile",
        "leopard",
        "leopard"
    ],
    "question" : "leopard",
    "date": "1502274893"
}
```
Both models are stored in the MongoDB.

The Captcha Model 'imgs' parameter, are UUIDs generated to set 'random' names to images. The server stores into MongoDB the relation between the 'random' name of each image and the real path of the image:
```json
{
    "captchaid" : "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
    "real" : "leopard/image_0092.jpg",
    "fake" : "1b838c46-b784-471e-b143-48be058c39a7.png"
}
```
When the server recieves a petition to get an image, recieves the petition with the fake image name, then, gets the real path of the image, gets it and serves the image content under the fake image name:
```
127.0.0.1:3025/image/1b838c46-b784-471e-b143-48be058c39a7.png
```

Captcha Model contains the captcha that server returns to the petition. And CaptchaSolution contains the solution of the captcha. Both have the same Id.


### 3.3 - Server validates captcha
When server recieves POST /answer, gets the answer, search for the CaptchaSolution based on the CaptchaId in the MongoDB, and then compares the answer 'selection' parameter with the CaptchaSolution.

If the selection is correct, returns 'true', if the selection is not correct, returns 'false'.


## 4 - Security

- If the captcha is resolved in less than 1 second, it's not valid.
- If the captcha is resolved in more than 1 minute, it's not valid.
- The images url, are UUIDs generated each time, in order to give different names for the images each time.
- The ip of requested captcha and answered captcha petitions must be the same.
- Each time a user fails answering the captcha, the server adds a counter to the IP and stores in to MongoDB. When the counter on that IP is greather than the value 'suspiciousIPCountLimit' defined in serverConfig.json, the IP is blocked for 'timeBan' seconds, also defined in serverConfig.json.
If before the counter exceeds the 'suspictiousIPCountLimit' the user answers correctly the captcha, the counter is deleted.
