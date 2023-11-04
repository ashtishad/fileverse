# Backend (Full Time) Assignment

APIs:

1. API for uploading a file:

```
  POST /upload

  Response: 
  {
    fileId: "",
    size: "",
    timestamp: ""
  }
```

2. API for serving a file:

```
  GET /file/:fileId

  Response: 
  <file content>
```


## Key points

* Handle large file 100MB+ (can change the API structure also)
* Documentation
* Folder Structure

## Bonus points

* Explain how to write unit tests
* Explain how to write end 2 end tests
* Postman Collection for easily testing
* Use IPFS for storage
* Docker File or Deployment Strategy

## Technologies

* Go
