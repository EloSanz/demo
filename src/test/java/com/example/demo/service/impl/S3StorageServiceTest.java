package com.example.demo.service.impl;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import java.nio.charset.StandardCharsets;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import software.amazon.awssdk.core.ResponseBytes;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.model.DeleteObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectRequest;
import software.amazon.awssdk.services.s3.model.GetObjectResponse;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Request;
import software.amazon.awssdk.services.s3.model.ListObjectsV2Response;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.model.PutObjectResponse;
import software.amazon.awssdk.services.s3.model.S3Object;

@ExtendWith(MockitoExtension.class)
class S3StorageServiceTest {

    @Mock private S3Client s3Client;

    @InjectMocks private S3StorageServiceImpl storageService;

    @Test
    void givenFileContent_whenUploadFile_shouldReturnSuccessMessage() {
        // Given
        String fileName = "test.txt";
        byte[] content = "hello world".getBytes(StandardCharsets.UTF_8);
        when(s3Client.putObject(any(PutObjectRequest.class), any(RequestBody.class)))
                .thenReturn(PutObjectResponse.builder().build());

        // When
        String result = storageService.uploadFile(fileName, content);

        // Then
        assertThat(result).isEqualTo("File uploaded: " + fileName);
        verify(s3Client).putObject(any(PutObjectRequest.class), any(RequestBody.class));
    }

    @Test
    void givenExistingFile_whenDownloadFile_shouldReturnFileBytes() {
        // Given
        String fileName = "test.txt";
        byte[] expectedContent = "hello world".getBytes(StandardCharsets.UTF_8);

        @SuppressWarnings("unchecked")
        ResponseBytes<GetObjectResponse> responseBytes =
                (ResponseBytes<GetObjectResponse>)
                        ResponseBytes.fromByteArray(
                                GetObjectResponse.builder().build(), expectedContent);

        when(s3Client.getObjectAsBytes(any(GetObjectRequest.class))).thenReturn(responseBytes);

        // When
        byte[] result = storageService.downloadFile(fileName);

        // Then
        assertThat(result).isEqualTo(expectedContent);
        verify(s3Client).getObjectAsBytes(any(GetObjectRequest.class));
    }

    @Test
    void givenExistingFile_whenDeleteFile_shouldInvokeS3Delete() {
        // Given
        String fileName = "to-delete.txt";

        // When
        storageService.deleteFile(fileName);

        // Then
        verify(s3Client).deleteObject(any(DeleteObjectRequest.class));
    }

    @Test
    void whenListFiles_shouldReturnFileNames() {
        // Given
        S3Object obj1 = S3Object.builder().key("file1.txt").build();
        S3Object obj2 = S3Object.builder().key("file2.jpg").build();
        ListObjectsV2Response listResponse =
                ListObjectsV2Response.builder().contents(obj1, obj2).build();
        when(s3Client.listObjectsV2(any(ListObjectsV2Request.class))).thenReturn(listResponse);

        // When
        List<String> files = storageService.listFiles();

        // Then
        assertThat(files).containsExactly("file1.txt", "file2.jpg");
        verify(s3Client).listObjectsV2(any(ListObjectsV2Request.class));
    }

    @Test
    void whenListFiles_andBucketIsEmpty_shouldReturnEmptyList() {
        // Given
        ListObjectsV2Response emptyResponse =
                ListObjectsV2Response.builder().contents(List.of()).build();
        when(s3Client.listObjectsV2(any(ListObjectsV2Request.class))).thenReturn(emptyResponse);

        // When
        List<String> files = storageService.listFiles();

        // Then
        assertThat(files).isEmpty();
    }
}
