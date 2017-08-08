# goCaptcha
captcha server, with own datasets, to train own machine learning AI


### How to use?

1. Get the captcha:
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
    "date": ""
}
```

2. User selects the images that fit in the 'question' parameter
(in this case, 'leopard')

3. Post the answer. The answer contains the CaptchaId, and an array with the selected images
```
POST/ 127.0.0.1:3025/answer
```
Post example:
```json
{
	"captchaid": "881c6083-0643-4d1c-9987-f8cc5bb9d5b1",
	"selection": [0,0,0,0,1,1]
}
```
Server response:
```
true
```

### How this works?

###### Server reads dataset
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


###### Server generates captcha
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
    "question" : "leopard"
}
```
Both models are stored in the MongoDB.

Captcha Model contains the captcha that server returns to the petition. And CaptchaSolution contains the solution of the captcha. Both have the same Id.


###### Server validates captcha
When server recieves POST /answer, gets the answer, search for the CaptchaSolution based on the CaptchaId in the MongoDB, and then compares the answer 'selection' parameter with the CaptchaSolution.

If the selection is correct, returns 'true', if the selection is not correct, returns 'false'.
