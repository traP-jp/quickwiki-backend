import re
from plantuml import PlantUML

def parse_sql_schema(file_path):
    with open(file_path, 'r') as file:
        sql_content = file.read()

    tables = re.findall(r'CREATE TABLE (\w+) \((.*?)\);', sql_content, re.S)
    schema = {}
    foreign_keys = []

    for table, columns in tables:
        column_list = re.findall(r'(\w+ .*?)(?:,|\n)', columns)
        schema[table] = column_list

        for column in column_list:
            if 'FOREIGN KEY' in column:
                fk_match = re.search(r'FOREIGN KEY \((\w+)\) REFERENCES (\w+)\((\w+)\)', column)
                if fk_match:
                    foreign_keys.append((table, fk_match.group(1), fk_match.group(2), fk_match.group(3)))

    return schema, foreign_keys

def generate_uml(schema, foreign_keys, output_path):
    uml_code = "@startuml\n\n"
    uml_code += "!theme carbon-gray\n"
    uml_code += "top to bottom direction\n"
    uml_code += "skinparam linetype ortho\n\n"
    for table, columns in schema.items():
        uml_code += f"class {table} {{\n"
        for column in columns:
            col_split = column.split(' ')
            if col_split[0] == 'FOREIGN' or col_split[0] == 'UNIQUE':
                continue
            uml_code += f"  {col_split[0]} : <color:#aaaaaa>{col_split[1].lower()}</color>\n"
        uml_code += "}\n"

    for table, fk_column, ref_table, ref_column in foreign_keys:
        uml_code += f"{table} -[#595959,plain]-^ {ref_table} : {fk_column} -> {ref_column}\n"

    uml_code += "@enduml"

    with open(output_path, 'w') as file:
        file.write(uml_code)

    plantuml = PlantUML(url='http://www.plantuml.com/plantuml/img/')
    plantuml.processes_file(output_path)

def generate_markdown(schema, output_path):
    with open(output_path, 'w') as file:
        file.write("# Database Schema\n\n")
        file.write('## TOC\n')
        for table in schema.keys():
            file.write(f'- [{table}](#{table.lower()})\n')
        file.write("## Tables\n\n")
        for table, columns in schema.items():
            file.write(f"### {table}\n\n")
            file.write("| Column Name | Data Type | Constraints |\n")
            file.write("|-------------|-----------|-------------|\n")

            foreign_key = []

            for column in columns:
                col_split = column.split(' ', 2)
                col_name = col_split[0]
                col_type = col_split[1]
                col_constraints = col_split[2] if len(col_split) == 3 else ''
                if col_name == 'FOREIGN':
                    _, _, key_name, _, reference = column.split(' ', 4)
                    foreign_key.append({'key_name': key_name, 'reference': reference})
                    continue
                if col_name == 'UNIQUE':
                    continue
                file.write(f"| {col_name} | {col_type.lower()} | {col_constraints.lower()} |\n")
            file.write("\n")
            if len(foreign_key) > 0:
                file.write("##### Foreign Keys\n")
                file.write("| Key Name | Reference |\n")
                file.write("|----------|-----------|\n")
                for fk in foreign_key:
                    file.write(f"| {fk['key_name']} | {fk['reference']} |\n")

        file.write("\n")
        file.write("## Diagram\n")
        file.write("![dbschema.png](dbschema.png)")

if __name__ == "__main__":
    schema, foreign_keys = parse_sql_schema('../schema.sql')
    generate_uml(schema, foreign_keys, 'dbschema.uml')
    generate_markdown(schema, 'dbschema.md')

