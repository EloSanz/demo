package com.example.demo.aspect;

import com.example.demo.exception.DomainException;
import java.time.Duration;
import java.time.Instant;
import java.util.Arrays;
import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.aspectj.lang.reflect.MethodSignature;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;

/**
 * Aspect for automatic logging of methods annotated with @Loggable. Replaces manual logging in
 * controllers/services.
 */
@Aspect
@Component
public class LoggingAspect {

    @Around("@within(Loggable) || @annotation(Loggable)")
    public Object logMethodExecution(ProceedingJoinPoint joinPoint) throws Throwable {
        MethodSignature signature = (MethodSignature) joinPoint.getSignature();
        Logger targetLogger = LoggerFactory.getLogger(signature.getDeclaringType());

        String className = signature.getDeclaringType().getSimpleName();
        String methodName = signature.getName();
        Loggable config = getLoggableConfig(signature);

        logMethodEntry(targetLogger, className, methodName, joinPoint.getArgs(), config);

        Instant start = Instant.now();
        try {
            Object result = joinPoint.proceed();
            long timeElapsed = Duration.between(start, Instant.now()).toMillis();

            logMethodExit(targetLogger, className, methodName, result, timeElapsed, config);
            return result;

        } catch (Throwable ex) {
            if (ex instanceof DomainException || ex instanceof IllegalArgumentException) {
                targetLogger.warn(
                        "<-- {}.{}() - FAILED: {}", className, methodName, ex.getMessage());
            } else {
                targetLogger.error(
                        "<-- {}.{}() - FAILED: {}", className, methodName, ex.getMessage(), ex);
            }
            throw ex;
        }
    }

    /** Extracts the @Loggable annotation either from the specific method or the class. */
    private Loggable getLoggableConfig(MethodSignature signature) {
        Loggable methodAnnotation = signature.getMethod().getAnnotation(Loggable.class);
        if (methodAnnotation != null) {
            return methodAnnotation;
        }

        Class<?> declaringType = signature.getDeclaringType();
        return declaringType.getAnnotation(Loggable.class);
    }

    private void logMethodEntry(
            Logger logger, String className, String methodName, Object[] args, Loggable config) {
        boolean shouldLogArgs = config != null && config.logArgs();

        if (shouldLogArgs && args != null && args.length > 0) {
            logger.info("--> {}.{}() - Args: {}", className, methodName, Arrays.toString(args));
        } else {
            logger.info("--> {}.{}()", className, methodName);
        }
    }

    private void logMethodExit(
            Logger logger,
            String className,
            String methodName,
            Object result,
            long timeElapsed,
            Loggable config) {
        boolean shouldLogResult = config != null && config.logResult();

        if (shouldLogResult) {
            logger.info(
                    "<-- {}.{}() - Result: {} - Took: {}ms",
                    className,
                    methodName,
                    result,
                    timeElapsed);
        } else {
            logger.info("<-- {}.{}() - Completed - Took: {}ms", className, methodName, timeElapsed);
        }
    }
}
