package com.example.demo.filter;

import static org.assertj.core.api.Assertions.assertThat;
import static org.springframework.boot.test.context.SpringBootTest.WebEnvironment.RANDOM_PORT;

import com.example.demo.dto.users.UserRequest;
import com.example.demo.util.UserTestUtil;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.boot.test.system.CapturedOutput;
import org.springframework.boot.test.system.OutputCaptureExtension;
import org.springframework.boot.test.web.client.TestRestTemplate;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;

@SpringBootTest(webEnvironment = RANDOM_PORT)
@ExtendWith(OutputCaptureExtension.class)
public class RequestResponseLoggingFilterTest {

    @Autowired
    private TestRestTemplate restTemplate;

    @Test
    public void shouldLogRequestAndResponseWithMaskedSensitiveData(CapturedOutput output) {
        // Given
        UserRequest request = UserTestUtil.createDefaultUserRequest();
        String sensitiveEmail = request.getEmail();

        // When
        ResponseEntity<String> response = restTemplate
                .postForEntity("/api/users", request, String.class);

        // Then
        assertThat(response.getStatusCode()).isIn(HttpStatus.CREATED, HttpStatus.BAD_REQUEST);

        // Verify request log contains masked email
        assertThat(output.getOut())
                .contains("\"email\":\"*******\"");

        // Verify output does NOT contain the actual email in the body
        // Note: The controller response might contain the email if it returns the
        // created user.
        // The filter logs the *request* and *response*.
        // If the response body contains the email (e.g. created user), it should also
        // be masked in the log.
        // However, standard REST APIs return the created resource.

        // Check logs for masked values
        assertThat(output.getOut())
                .doesNotContain("\"email\":\"" + sensitiveEmail + "\"");
    }
}
