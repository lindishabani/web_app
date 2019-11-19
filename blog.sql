-- phpMyAdmin SQL Dump
-- version 4.8.4
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Nov 19, 2019 at 10:32 PM
-- Server version: 10.1.37-MariaDB
-- PHP Version: 7.3.1

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `blog`
--

-- --------------------------------------------------------

--
-- Table structure for table `category`
--

CREATE TABLE `category` (
  `id` int(11) NOT NULL,
  `category` varchar(255) NOT NULL,
  `date` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `category`
--

INSERT INTO `category` (`id`, `category`, `date`) VALUES
(12, 'Programming language', '2019-11-13 22:18:23'),
(13, 'Design', '2019-11-13 22:19:50');

-- --------------------------------------------------------

--
-- Table structure for table `comments`
--

CREATE TABLE `comments` (
  `id` int(11) NOT NULL,
  `comment` varchar(255) NOT NULL,
  `date` varchar(255) NOT NULL,
  `post_related` int(11) NOT NULL,
  `comment_author` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `comments`
--

INSERT INTO `comments` (`id`, `comment`, `date`, `post_related`, `comment_author`) VALUES
(2, 'Photoshop', '2019-11-13 23:36:01', 4, 1),
(3, 'Golang', '2019-11-14 00:21:18', 5, 1),
(4, 'asdsad', '2019-11-14 17:30:22', 5, 2);

-- --------------------------------------------------------

--
-- Table structure for table `posts`
--

CREATE TABLE `posts` (
  `id` int(11) NOT NULL,
  `title` varchar(255) NOT NULL,
  `description` varchar(10000) NOT NULL,
  `category_id` int(11) NOT NULL,
  `date` varchar(255) NOT NULL,
  `author_id` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `posts`
--

INSERT INTO `posts` (`id`, `title`, `description`, `category_id`, `date`, `author_id`) VALUES
(2, 'Golang', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas viverra nisi lacus. Mauris quis velit sed massa pharetra sodales. Aliquam consequat varius metus, sit amet auctor sem rhoncus sed. Curabitur sodales ante eget velit pulvinar, in feugiat lib', 13, '2019-11-14 16:06:20', 1),
(3, 'Python', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas viverra nisi lacus. Mauris quis velit sed massa pharetra sodales. Aliquam consequat varius metus, sit amet auctor sem rhoncus sed. Curabitur sodales ante eget velit pulvinar, in feugiat lib', 12, '2019-11-13 22:30:39', 1),
(4, 'Photoshop', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas viverra nisi lacus. Mauris quis velit sed massa pharetra sodales. Aliquam consequat varius metus, sit amet auctor sem rhoncus sed. Curabitur sodales ante eget velit pulvinar, in feugiat lib', 13, '2019-11-13 22:38:16', 1),
(5, 'Ilustrator', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Maecenas viverra nisi lacus. Mauris quis velit sed massa pharetra sodales. Aliquam consequat varius metus, sit amet auctor sem rhoncus sed. Curabitur sodales ante eget velit pulvinar, in feugiat lib', 13, '2019-11-13 22:46:46', 1);

-- --------------------------------------------------------

--
-- Table structure for table `users`
--

CREATE TABLE `users` (
  `id` int(11) NOT NULL,
  `name` varchar(255) NOT NULL,
  `surname` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` int(11) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

--
-- Dumping data for table `users`
--

INSERT INTO `users` (`id`, `name`, `surname`, `email`, `password`, `role`) VALUES
(1, 'lind', 'shabani', 'lindshabani97@hotmail.com', '$2a$14$DiNR0X5vz1D.2GgYwlVN/u5YsDHe58oWyVFPlEkVC/Gg3hSb4DFLi', 1),
(2, 'lindi', 'shabani', 'lindi@shabani.com', '$2a$14$qjKLLVyVNNGv4N5bFlPYq.w1aOkIIAWWsRwDbJrf1vQdjz9q0Aqbi', 0);

--
-- Indexes for dumped tables
--

--
-- Indexes for table `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `comments`
--
ALTER TABLE `comments`
  ADD PRIMARY KEY (`id`),
  ADD KEY `comment_author` (`comment_author`),
  ADD KEY `post_related` (`post_related`);

--
-- Indexes for table `posts`
--
ALTER TABLE `posts`
  ADD PRIMARY KEY (`id`),
  ADD KEY `category_id` (`category_id`),
  ADD KEY `posts_ibfk_1` (`author_id`);

--
-- Indexes for table `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `category`
--
ALTER TABLE `category`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- AUTO_INCREMENT for table `comments`
--
ALTER TABLE `comments`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=5;

--
-- AUTO_INCREMENT for table `posts`
--
ALTER TABLE `posts`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=6;

--
-- AUTO_INCREMENT for table `users`
--
ALTER TABLE `users`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=3;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `comments`
--
ALTER TABLE `comments`
  ADD CONSTRAINT `comments_ibfk_1` FOREIGN KEY (`comment_author`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `comments_ibfk_2` FOREIGN KEY (`post_related`) REFERENCES `posts` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;

--
-- Constraints for table `posts`
--
ALTER TABLE `posts`
  ADD CONSTRAINT `posts_ibfk_1` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE,
  ADD CONSTRAINT `posts_ibfk_2` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE ON UPDATE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
