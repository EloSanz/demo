package com.example.demo.aspect;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mockStatic;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.example.demo.exception.DomainException;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.reflect.MethodSignature;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.MockedStatic;
import org.mockito.junit.jupiter.MockitoExtension;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

@ExtendWith(MockitoExtension.class)
class LoggingAspectTest {

    private LoggingAspect loggingAspect;

    @Mock private ProceedingJoinPoint joinPoint;

    @Mock private MethodSignature signature;

    @Mock private Logger logger;

    private MockedStatic<LoggerFactory> loggerFactoryMockedStatic;

    @BeforeEach
    @SuppressWarnings({"unchecked", "rawtypes"})
    void setUp() throws NoSuchMethodException {
        loggingAspect = new LoggingAspect();

        loggerFactoryMockedStatic = mockStatic(LoggerFactory.class);
        loggerFactoryMockedStatic
                .when(() -> LoggerFactory.getLogger(any(Class.class)))
                .thenReturn(logger);

        when(joinPoint.getSignature()).thenReturn(signature);
        when(signature.getDeclaringType()).thenReturn((Class) TestTargetClass.class);
        when(signature.getName()).thenReturn("testMethod");
        when(signature.getMethod()).thenReturn(TestTargetClass.class.getMethod("testMethod"));
    }

    @AfterEach
    void tearDown() {
        if (loggerFactoryMockedStatic != null) {
            loggerFactoryMockedStatic.close();
        }
    }

    @Test
    void logMethodExecution_Success_WithArgsAndResult() throws Throwable {
        when(joinPoint.getArgs()).thenReturn(new Object[] {"arg1"});
        when(joinPoint.proceed()).thenReturn("result");

        Object result = loggingAspect.logMethodExecution(joinPoint);

        assertEquals("result", result);
        verify(joinPoint, times(1)).proceed();

        verify(logger)
                .info(
                        eq("--> {}.{}() - Args: {}"),
                        eq("TestTargetClass"),
                        eq("testMethod"),
                        eq("[arg1]"));
        verify(logger)
                .info(
                        eq("<-- {}.{}() - Result: {} - Took: {}ms"),
                        eq("TestTargetClass"),
                        eq("testMethod"),
                        eq("result"),
                        any(Long.class));
    }

    @Test
    void logMethodExecution_Success_WithoutArgsAndResult() throws Throwable {
        when(signature.getMethod()).thenReturn(TestTargetClassNoLog.class.getMethod("testMethod"));
        when(joinPoint.getArgs()).thenReturn(new Object[] {"arg1"});
        when(joinPoint.proceed()).thenReturn("result");

        Object result = loggingAspect.logMethodExecution(joinPoint);

        assertEquals("result", result);
        verify(logger).info("--> {}.{}()", "TestTargetClass", "testMethod");
        verify(logger)
                .info(
                        eq("<-- {}.{}() - Completed - Took: {}ms"),
                        eq("TestTargetClass"),
                        eq("testMethod"),
                        any(Long.class));
    }

    @Test
    void logMethodExecution_ThrowsDomainException_LogsWarn() throws Throwable {
        class TestDomainException extends DomainException {
            public TestDomainException(String message) {
                super(message);
            }
        }

        TestDomainException exception = new TestDomainException("Domain error");
        when(joinPoint.proceed()).thenThrow(exception);

        assertThrows(TestDomainException.class, () -> loggingAspect.logMethodExecution(joinPoint));

        verify(logger)
                .warn("<-- {}.{}() - FAILED: {}", "TestTargetClass", "testMethod", "Domain error");
    }

    @Test
    void logMethodExecution_ThrowsUnexpectedException_LogsError() throws Throwable {
        RuntimeException exception = new RuntimeException("Unexpected error");
        when(joinPoint.proceed()).thenThrow(exception);

        assertThrows(RuntimeException.class, () -> loggingAspect.logMethodExecution(joinPoint));

        verify(logger)
                .error(
                        "<-- {}.{}() - FAILED: {}",
                        "TestTargetClass",
                        "testMethod",
                        "Unexpected error",
                        exception);
    }

    static class TestTargetClass {
        @Loggable(logArgs = true, logResult = true)
        public String testMethod() {
            return "test";
        }
    }

    static class TestTargetClassNoLog {
        @Loggable(logArgs = false, logResult = false)
        public String testMethod() {
            return "test";
        }
    }
}
