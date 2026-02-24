package com.example.demo.service.impl;

import com.example.demo.exception.storage.StorageException;
import com.example.demo.service.StorageService;
import java.util.List;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import software.amazon.awssdk.core.ResponseBytes;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.DeleteObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectResponse;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Request;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Response;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.model.S3Object;

@Service
@RequiredArgsConstructor
public class S3StorageServiceImpl implements StorageService {

    private final S3Client s3Client;

    @Value("${aws.s3.bucket-name:myawsbucketelito}")
    private String bucketName;

    @Override
    public String uploadFile(String fileName, byte[] content) {
        try {
            PutObjectRequest putOb =
                    PutObjectRequest.builder().bucket(bucketName).key(fileName).build();
            s3Client.putObject(putOb, RequestBody.fromBytes(content));
            return "File uploaded: " + fileName;
        } catch (Exception e) {
            throw new StorageException("Failed to upload file: " + fileName, e);
        }
    }

    @Override
    public byte[] downloadFile(String fileName) {
        try {
            GetObjectRequest getOb =
                    GetObjectRequest.builder().bucket(bucketName).key(fileName).build();
            ResponseBytes<GetObjectResponse> objectBytes = s3Client.getObjectAsBytes(getOb);
            return objectBytes.asByteArray();
        } catch (Exception e) {
            throw new StorageException("Failed to download file: " + fileName, e);
        }
    }

    @Override
    public void deleteFile(String fileName) {
        try {
            DeleteObjectRequest deleteOb =
                    DeleteObjectRequest.builder().bucket(bucketName).key(fileName).build();
            s3Client.deleteObject(deleteOb);
        } catch (Exception e) {
            throw new StorageException("Failed to delete file: " + fileName, e);
        }
    }

    @Override
    public List<String> listFiles() {
        try {
            ListObjectsV2Request listOb = ListObjectsV2Request.builder().bucket(bucketName).build();
            ListObjectsV2Response response = s3Client.listObjectsV2(listOb);
            return response.contents().stream().map(S3Object::key).toList();
        } catch (Exception e) {
            throw new StorageException("Failed to list files in bucket", e);
        }
    }
}
