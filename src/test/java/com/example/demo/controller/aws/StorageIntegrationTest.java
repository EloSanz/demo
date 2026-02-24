package com.example.demo.controller.aws;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.example.demo.BaseIntegrationTest;
import com.example.demo.service.StorageService;
import java.nio.charset.StandardCharsets;
import java.util.List;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.core.ParameterizedTypeReference;
import org.springframework.core.io.ByteArrayResource;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.util.LinkedMultiValueMap;
import org.springframework.util.MultiValueMap;

class StorageIntegrationTest extends BaseIntegrationTest {

    @Autowired private TestRestTemplate restTemplate;

    @MockitoBean private StorageService storageService;

    @Test
    void givenFile_whenUpload_shouldReturn200WithMessage() {
        // Given
        String fileName = "hello.txt";
        byte[] fileContent = "hello world".getBytes(StandardCharsets.UTF_8);
        when(storageService.uploadFile(eq(fileName), any(byte[].class)))
                .thenReturn("File uploaded: " + fileName);

        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.MULTIPART_FORM_DATA);

        MultiValueMap<String, Object> body = new LinkedMultiValueMap<>();
        ByteArrayResource resource =
                new ByteArrayResource(fileContent) {
                    @Override
                    public String getFilename() {
                        return fileName;
                    }
                };
        body.add("file", resource);

        // When
        ResponseEntity<String> response =
                restTemplate.postForEntity(
                        "/api/storage/upload", new HttpEntity<>(body, headers), String.class);

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(response.getBody()).isEqualTo("File uploaded: " + fileName);
        verify(storageService).uploadFile(eq(fileName), any(byte[].class));
    }

    @Test
    void whenListFiles_shouldReturnFileNames() {
        // Given
        when(storageService.listFiles()).thenReturn(List.of("file1.txt", "file2.jpg"));

        // When
        ResponseEntity<List<String>> response =
                restTemplate.exchange(
                        "/api/storage/files",
                        HttpMethod.GET,
                        null,
                        new ParameterizedTypeReference<>() {});

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(response.getBody()).containsExactly("file1.txt", "file2.jpg");
    }

    @Test
    void whenListFiles_andBucketIsEmpty_shouldReturnEmptyList() {
        // Given
        when(storageService.listFiles()).thenReturn(List.of());

        // When
        ResponseEntity<List<String>> response =
                restTemplate.exchange(
                        "/api/storage/files",
                        HttpMethod.GET,
                        null,
                        new ParameterizedTypeReference<>() {});

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(response.getBody()).isEmpty();
    }

    @Test
    void givenExistingFile_whenDownload_shouldReturnFileBytes() {
        // Given
        String fileName = "hello.txt";
        byte[] expectedContent = "hello world".getBytes(StandardCharsets.UTF_8);
        when(storageService.downloadFile(fileName)).thenReturn(expectedContent);

        // When
        ResponseEntity<byte[]> response =
                restTemplate.getForEntity("/api/storage/files/" + fileName, byte[].class);

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.OK);
        assertThat(response.getBody()).isEqualTo(expectedContent);
        assertThat(response.getHeaders().getContentType())
                .isEqualTo(MediaType.APPLICATION_OCTET_STREAM);
        verify(storageService).downloadFile(fileName);
    }

    @Test
    void givenExistingFile_whenDelete_shouldReturn204() {
        // Given
        String fileName = "to-delete.txt";
        doNothing().when(storageService).deleteFile(fileName);

        // When
        ResponseEntity<Void> response =
                restTemplate.exchange(
                        "/api/storage/files/" + fileName, HttpMethod.DELETE, null, Void.class);

        // Then
        assertThat(response.getStatusCode()).isEqualTo(HttpStatus.NO_CONTENT);
        verify(storageService).deleteFile(fileName);
    }
}
