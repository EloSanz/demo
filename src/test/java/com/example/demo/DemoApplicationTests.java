package com.example.demo;

import org.junit.jupiter.api.Test;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import software.amazon.awssdk.services.s3.S3Client;

@SpringBootTest
class DemoApplicationTests {

    @MockitoBean private S3Client s3Client;

    @Test
    void contextLoads() {}
}
