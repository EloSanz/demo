package com.example.demo;

import com.example.demo.repository.UserRepository;
import org.junit.jupiter.api.BeforeEach;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.web.server.LocalServerPort;
import org.springframework.test.context.ActiveProfiles;

/**
 * Base class for integration tests. Uses H2 in-memory database for fast, dependency-free execution.
 */
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@ActiveProfiles("default") // Uses H2 from default profile
public abstract class BaseIntegrationTest {

    @LocalServerPort protected int port;

    @Autowired protected UserRepository userRepository;

    @BeforeEach
    void setUp() {
        userRepository.deleteAll();
    }
}
