package com.example.demo.controller.users;

import com.example.demo.dto.users.UserResponse;
import com.example.demo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.responses.ApiResponses;
import java.util.List;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@Slf4j
public class ExternalSyncEndpoint extends BaseUserController {

    public ExternalSyncEndpoint(UserService userService) {
        super(userService);
    }

    @Operation(
            summary = "Fetch external users",
            description = "Retrieves users from the external JSONPlaceholder API.")
    @ApiResponse(responseCode = "200", description = "External users retrieved successfully")
    @GetMapping("/external")
    public ResponseEntity<List<UserResponse>> fetchExternalUsers() {
        log.info("GET /api/users/external - Fetching users from external API");
        return ResponseEntity.ok(userService.fetchUsersFromExternalApi());
    }

    @Operation(
            summary = "Sync external user",
            description =
                    "Fetches a user from the external API and saves them to the local database.")
    @ApiResponses(
            value = {
                @ApiResponse(responseCode = "201", description = "User synced and created locally"),
                @ApiResponse(responseCode = "404", description = "User not found in external API")
            })
    @PostMapping("/sync/{externalUserId}")
    public ResponseEntity<UserResponse> syncUserFromExternal(
            @Parameter(description = "ID of the user in the external system") @PathVariable
                    Long externalUserId) {
        log.info("POST /api/users/sync/{} - Syncing user from external API", externalUserId);
        UserResponse syncedUser = userService.syncUserFromExternalApi(externalUserId);
        return ResponseEntity.status(HttpStatus.CREATED).body(syncedUser);
    }
}
