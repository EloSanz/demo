package com.example.demo;

import org.springframework.boot.test.context.SpringBootTest;

@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
public abstract class AbstractIntegrationTest {
    // H2 will be used by default due to 'com.h2database:h2' in dependencies
    // and missing postgres driver (or strictly prioritized H2 context)
}
