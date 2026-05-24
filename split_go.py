import re
import os

def split_methods(source_file, dest_dir, struct_name, method_groups, common_imports, is_controller=False):
    with open(source_file, 'r') as f:
        content = f.read()

    # Find the start of the struct methods
    # We look for "func (receiver *StructName) MethodName("
    method_pattern = r'(?m)^func \(\w+ \*' + struct_name + r'\) (\w+)\(.*?$'
    
    methods = {}
    
    # We need to extract the exact code blocks for each method
    # Since Go format puts func starting at line start, we can match lines starting with "func "
    lines = content.split('\n')
    
    current_method = None
    method_code = []
    
    header = []
    in_header = True
    
    for line in lines:
        if line.startswith(f'func ('):
            match = re.search(r'func \(\w+ \*' + struct_name + r'\) (\w+)\(', line)
            if match:
                if current_method:
                    methods[current_method] = '\n'.join(method_code)
                current_method = match.group(1)
                method_code = [line]
                in_header = False
                continue
        
        if in_header:
            header.append(line)
        elif current_method:
            method_code.append(line)
            
    if current_method:
        methods[current_method] = '\n'.join(method_code)
        
    # Write the base file with struct def and New... function
    base_file = os.path.join(dest_dir, 'base.go')
    # we just rename header to what we need
    # but some non-method funcs like 'success', 'isBlank' might be in method_code if not careful
    pass # we'll use a safer approach for the base file

def run_split():
    # Since regex can be tricky with nested brackets, let's do a simpler text-based chunking.
    
    def extract_chunks(filepath, struct_name):
        with open(filepath, 'r') as f:
            lines = f.readlines()
            
        chunks = {}
        base_lines = []
        
        current_method = None
        current_chunk = []
        
        for line in lines:
            if line.startswith('func (') and f'*{struct_name}' in line:
                if current_method:
                    chunks[current_method] = ''.join(current_chunk)
                match = re.search(r'\) (\w+)\(', line)
                if match:
                    current_method = match.group(1)
                else:
                    current_method = "UNKNOWN"
                current_chunk = [line]
            else:
                if current_method is None:
                    base_lines.append(line)
                else:
                    current_chunk.append(line)
                    
        if current_method:
            chunks[current_method] = ''.join(current_chunk)
            
        return base_lines, chunks

    svc_base, svc_chunks = extract_chunks('apps/fake-uit-server/internal/service/usecase.go', 'Service')
    ctrl_base, ctrl_chunks = extract_chunks('apps/fake-uit-server/internal/controller/uit_controller.go', 'UitController')
    
    groups = {
        'auth': ['Login', 'GetStudentByIDForAuth'],
        'room': ['GetRoomsAvailability'],
        'course': ['GetDeadlines', 'GetMaterials', 'GetAssignments', 'SubmitAssignment'],
        'ctsv': ['ConfirmLetter', 'BankLoans', 'TrainingPointConfirm', 'LanguageCertificate'],
    }
    
    def write_group(dest_file, pkg, chunks, group_name):
        methods_to_write = groups.get(group_name, [])
        code = f"package {pkg}\n\n"
        if pkg == "controller":
            code += 'import (\n\t"net/http"\n\t"github.com/gin-gonic/gin"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"\n)\n\n'
        else:
            code += 'import (\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"\n\t"github.com/google/uuid"\n\t"time"\n\t"net/http"\n\t"strings"\n)\n\n'
            
        wrote_any = False
        for m in methods_to_write:
            if m in chunks:
                code += chunks[m] + "\n"
                del chunks[m]
                wrote_any = True
                
        if wrote_any:
            with open(dest_file, 'w') as f:
                f.write(code)

    for g in groups.keys():
        write_group(f'apps/fake-uit-server/internal/service/{g}_service.go', 'service', svc_chunks, g)
        write_group(f'apps/fake-uit-server/internal/controller/{g}_controller.go', 'controller', ctrl_chunks, g)
        
    # The rest goes to student
    student_svc = "package service\n\n"
    student_svc += 'import (\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"\n\t"github.com/google/uuid"\n\t"time"\n\t"strings"\n\t"net/http"\n)\n\n'
    for m, c in list(svc_chunks.items()):
        student_svc += c + "\n"
        
    with open('apps/fake-uit-server/internal/service/student_service.go', 'w') as f:
        f.write(student_svc)

    student_ctrl = "package controller\n\n"
    student_ctrl += 'import (\n\t"net/http"\n\t"github.com/gin-gonic/gin"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/dto"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/apperror"\n\t"github.com/giakiet05/uit-hub/apps/fake-uit-server/internal/middleware"\n)\n\n'
    for m, c in list(ctrl_chunks.items()):
        student_ctrl += c + "\n"
        
    with open('apps/fake-uit-server/internal/controller/student_controller.go', 'w') as f:
        f.write(student_ctrl)

    # Write base files
    with open('apps/fake-uit-server/internal/service/service.go', 'w') as f:
        f.write(''.join(svc_base))
        
    with open('apps/fake-uit-server/internal/controller/controller.go', 'w') as f:
        f.write(''.join(ctrl_base))
        
    os.remove('apps/fake-uit-server/internal/service/usecase.go')
    os.remove('apps/fake-uit-server/internal/controller/uit_controller.go')
    
run_split()
