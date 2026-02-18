package com.example.demo.filter;

import jakarta.servlet.FilterChain;
import jakarta.servlet.ServletException;
import jakarta.servlet.http.HttpServletRequest;
import jakarta.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.util.Enumeration;
import java.util.HashMap;
import java.util.Map;
import java.util.Set;
import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Component;
import org.springframework.web.filter.OncePerRequestFilter;
import org.springframework.web.util.ContentCachingRequestWrapper;
import org.springframework.web.util.ContentCachingResponseWrapper;

@Component
@Slf4j
public class RequestResponseLoggingFilter extends OncePerRequestFilter {

    private static final Set<String> SENSITIVE_FIELDS = Set.of("email", "password", "token");

    @Override
    protected void doFilterInternal(
            HttpServletRequest request, HttpServletResponse response, FilterChain filterChain)
            throws ServletException, IOException {

        ContentCachingRequestWrapper requestWrapper = new ContentCachingRequestWrapper(request);
        ContentCachingResponseWrapper responseWrapper = new ContentCachingResponseWrapper(response);

        try {
            filterChain.doFilter(requestWrapper, responseWrapper);
        } finally {
            logRequest(requestWrapper);
            logResponse(responseWrapper);
            responseWrapper.copyBodyToResponse();
        }
    }

    private void logRequest(ContentCachingRequestWrapper request) {
        String method = request.getMethod();
        String uri = request.getRequestURI();
        String body = getBody(request.getContentAsByteArray());
        Map<String, String> headers = getHeaders(request);

        log.info(
                "Request: method={}, uri={}, headers={}, body={}",
                method,
                uri,
                headers,
                maskSensitiveData(body));
    }

    private void logResponse(ContentCachingResponseWrapper response) {
        int status = response.getStatus();
        String body = getBody(response.getContentAsByteArray());

        log.info("Response: status={}, body={}", status, maskSensitiveData(body));
    }

    private String getBody(byte[] content) {
        if (content.length == 0) {
            return "";
        }
        return new String(content, StandardCharsets.UTF_8);
    }

    private Map<String, String> getHeaders(HttpServletRequest request) {
        Map<String, String> headers = new HashMap<>();
        Enumeration<String> headerNames = request.getHeaderNames();
        while (headerNames.hasMoreElements()) {
            String headerName = headerNames.nextElement();
            headers.put(headerName, request.getHeader(headerName));
        }
        return headers;
    }

    // Simple regex-based masking for demonstration.
    // In production, use a JSON library like Jackson with a custom
    // serializer/module.
    private String maskSensitiveData(String json) {
        if (json == null || json.isEmpty()) {
            return json;
        }

        String maskedJson = json;
        for (String field : SENSITIVE_FIELDS) {
            maskedJson = maskedJson.replaceAll(
                    "(\"" + field + "\"\\s*:\\s*\")([^\"]+)(\")", "$1*******$3");
        }
        return maskedJson;
    }
}
