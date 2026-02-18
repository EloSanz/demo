package com.example.demo.config;

import com.example.demo.client.ExternalUserClient;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.reactive.function.client.WebClient;
import org.springframework.web.reactive.function.client.support.WebClientAdapter;
import org.springframework.web.service.invoker.HttpServiceProxyFactory;

/**
 * Configuration for HTTP Interface clients.
 *
 * <p>This creates the implementation of declarative HTTP interfaces using the base URL from
 * application.yml.
 */
@Configuration
public class HttpClientConfig {

    @Value("${external.api.jsonplaceholder.base-url}")
    private String jsonPlaceholderBaseUrl;

    /**
     * Creates a bean for the ExternalUserClient interface. Spring will generate the implementation
     * automatically.
     */
    @Bean
    public ExternalUserClient externalUserClient() {
        WebClient webClient = WebClient.builder().baseUrl(jsonPlaceholderBaseUrl).build();

        HttpServiceProxyFactory factory =
                HttpServiceProxyFactory.builderFor(WebClientAdapter.create(webClient)).build();

        return factory.createClient(ExternalUserClient.class);
    }
}
