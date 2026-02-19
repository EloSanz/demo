package com.example.demo.controller.users;

import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Slf4j
public class DeleteUserEndpoint extends BaseUserController {

    public DeleteUserEndpoint(UserService userService) {
        super(userService);
    }

    @Operation(summary = "Delete user", description = "Deletes a user by their ID.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "204", description = "User deleted successfully"),
                @ApiResponse(responseCode = "404", description = "User not found")
            })
    @DeleteMapping("/{id}")
    public ResponseEntity<Void> deleteUser(
            @Parameter(description = "ID of the user to delete") @PathVariable Long id) {
        userService.deleteUser(id);
        return ResponseEntity.noContent().build();
    }
}
