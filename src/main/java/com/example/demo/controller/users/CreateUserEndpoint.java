package com.example.demo.controller.users;

import com.example.demo.controller.dto.users.UserRequest;
import com.example.demo.controller.dto.users.UserResponse;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import jakarta.validation.Valid;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Slf4j
public class CreateUserEndpoint extends BaseUserController {

    public CreateUserEndpoint(UserService userService) {
        super(userService);
    }

    @Operation(
            summary = "Create a new user",
            description = "Creates a new user in the local database. Email must be unique.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "201", description = "User created successfully"),
                @ApiResponse(
                        responseCode = "400",
                        description = "Invalid input or email already exists")
            })
    @PostMapping
    public ResponseEntity<UserResponse> createUser(@Valid @RequestBody UserRequest request) {
        log.info("POST /api/users - Creating new user: {}", request.getEmail());
        UserResponse createdUser = userService.createUser(request);
        return ResponseEntity.status(HttpStatus.CREATED).body(createdUser);
    }
}
