package com.example.demo.config;

import io.swagger.v3.oas.models.OpenAPI;
import io.swagger.v3.oas.models.info.Contact;
import io.swagger.v3.oas.models.info.Info;
import io.swagger.v3.oas.models.info.License;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class OpenApiConfig {

    @Bean
    public OpenAPI customOpenAPI() {
        return new OpenAPI()
                .info(
                        new Info()
                                .title("Spring Boot 4 Template API")
                                .version("1.0.0")
                                .description(
                                        "Template project demonstrating clean architecture, HTTP Interfaces, and best practices.")
                                .contact(
                                        new Contact()
                                                .name("Tu Nombre")
                                                .email("tu@email.com")
                                                .url("https://tusitio.com"))
                                .license(
                                        new License()
                                                .name("MIT License")
                                                .url("https://opensource.org/licenses/MIT")));
    }
}
