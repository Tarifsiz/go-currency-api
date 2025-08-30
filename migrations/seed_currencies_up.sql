-- Insert currencies based on the provided data
-- Using a fixed UUID for created_by

INSERT INTO currencies (id, code, description, amount_display_format, html_encoded_symbol, factor, created_by) VALUES 
-- Original batch
('eb30cd07-76fe-4b24-b9f2-0274eabf304f', 'AED', 'United Arab Emirates Dirham', '###,###.##', '&#x62f;&#x2e;&#x625;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('4022f85b-c321-43e6-9e31-03f37c6de0b9', 'MAD', 'Moroccan Dirham', '###,###.##', '&#77;&#65;&#68;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('9d4d2f61-56d6-4893-a553-10d91edb7837', 'MUR', 'Mauritian Rupee', '###,###.##', '&#8360;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('3714940c-c4f1-4476-b7aa-15ffc56ba3e0', 'XCD', 'Eastern Caribbean Dollar', '###,###.##', '&#36;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('c0589eda-3b60-427f-a17b-183ec6e7b99d', 'CLP', 'Chilean Peso', '###,###', '&#36;', 1, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('4f3052d2-9f00-4cea-8a85-1e812cd52fe6', 'ZAR', 'South African Rand', '###,###.##', '&#82;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('03bd3b96-59ac-4329-9951-1ecffd3f7de7', 'SEK', 'Swedish Krona', '###,###.##', '&#107;&#114;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('e87d08bc-3fc1-4f52-91bd-1ed21153ada0', 'KES', 'Kenyan Shilling', '###,###.##', '&#75;&#83;&#104;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),

-- Additional batch from new data
('ac508807-d802-4cc2-bcda-27db132b7c06', 'CAD', 'Canadian Dollar', '###,###.##', '&#36;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('7c758828-db9e-4db6-87df-2cc1f54b7709', 'GBP', 'British Pound', '###,###.##', '&#163;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('ad03349f-0514-4bbc-98ac-2faa2dd41b46', 'OMR', 'Omani Rial', '###,###.###', '&#65020;', 1000, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('a57ed92d-d175-4e66-9090-3027ac9bb7ca', 'RON', 'Romanian Leu', '###,###.##', '&#108;&#101;&#105;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('8caccd02-dbdb-493a-ba77-3e8fe2b485ae', 'NOK', 'Norwegian Krone', '###,###.##', '&#107;&#114;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('220a4fa7-1432-495e-bc7b-4622376eae77', 'SAR', 'Saudi Riyal', '###,###.##', '&#65020;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('11514c9a-d6fe-4097-ad02-470a92d06a62', 'JPY', 'Japanese Yen', '###,###', '&#165;', 1, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('6fa52ee6-3e57-416b-aaa7-498309cafd14', 'DKK', 'Danish Krone', '###,###.##', '&#107;&#114;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('077b3d88-950c-41ae-b343-4e2d2634351d', 'HUF', 'Hungarian Forint', '###,###.##', '&#70;&#116;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('edbd44cd-c85b-4476-aa83-4fbe3f3c2182', 'IDR', 'Indonesian Rupiah', '###,###.##', '&#36;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24');

-- Add some other common major currencies for completeness
INSERT INTO currencies (code, description, amount_display_format, html_encoded_symbol, factor, created_by) VALUES 
('USD', 'United States Dollar', '###,###.##', '&#36;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('EUR', 'Euro', '###,###.##', '&#8364;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('CHF', 'Swiss Franc', '###,###.##', '&#67;&#72;&#70;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('CNY', 'Chinese Yuan', '###,###.##', '&#165;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('INR', 'Indian Rupee', '###,###.##', '&#8377;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('TRY', 'Turkish Lira', '###,###.##', '&#8378;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24'),
('AUD', 'Australian Dollar', '###,###.##', '&#36;', 100, '1609b0e1-30c4-402c-a76e-8f5b4d6cfc24');