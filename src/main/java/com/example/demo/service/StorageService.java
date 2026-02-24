package com.example.demo.service;

import java.util.List;

public interface StorageService {

    String uploadFile(String fileName, byte[] content);

    byte[] downloadFile(String fileName);

    void deleteFile(String fileName);

    List<String> listFiles();
}
