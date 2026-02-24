package com.example.demo.controller.aws;

import com.example.demo.aspect.Loggable;
import com.example.demo.service.StorageService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.tags.Tag;
import java.io.IOException;
import java.util.List;
import lombok.RequiredArgsConstructor;
import org.springframework.http.ContentDisposition;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.multipart.MultipartFile;

@RestController
@RequestMapping("/api/storage")
@RequiredArgsConstructor
@Tag(name = "Storage", description = "AWS S3 Storage Operations")
@Loggable(logArgs = true, logResult = false)
public class StorageController {

    private final StorageService storageService;

    @PostMapping(value = "/upload", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @Operation(summary = "Upload a file to Amazon S3")
    public ResponseEntity<String> uploadFile(@RequestParam("file") MultipartFile file) {
        try {
            String result = storageService.uploadFile(file.getOriginalFilename(), file.getBytes());
            return ResponseEntity.ok(result);
        } catch (IOException e) {
            return ResponseEntity.internalServerError()
                    .body("Error reading file: " + e.getMessage());
        }
    }

    @GetMapping("/files")
    @Operation(summary = "List all files in S3 bucket")
    public ResponseEntity<List<String>> listFiles() {
        List<String> files = storageService.listFiles();
        return ResponseEntity.ok(files);
    }

    @GetMapping("/files/{fileName}")
    @Operation(summary = "Download a file from Amazon S3")
    public ResponseEntity<byte[]> downloadFile(@PathVariable String fileName) {
        byte[] content = storageService.downloadFile(fileName);
        HttpHeaders headers = new HttpHeaders();
        headers.setContentType(MediaType.APPLICATION_OCTET_STREAM);
        headers.setContentDisposition(ContentDisposition.attachment().filename(fileName).build());
        return ResponseEntity.ok().headers(headers).body(content);
    }

    @DeleteMapping("/files/{fileName}")
    @Operation(summary = "Delete a file from Amazon S3")
    public ResponseEntity<Void> deleteFile(@PathVariable String fileName) {
        storageService.deleteFile(fileName);
        return ResponseEntity.noContent().build();
    }
}
