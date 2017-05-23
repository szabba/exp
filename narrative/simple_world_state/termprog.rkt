#lang rackjure

(require srfi/1)
(require charterm)

(provide program)

(provide (contract-out [program? predicate/c]
                       [run! (-> program? void?)]))


;;; Terminal-based program structure

(struct program (model update view exit?) #:transparent)

(define (run! prog)
  (with-charterm (run-loop! prog))
  (void))

       
(define (run-loop! prog)
  (draw-program! prog)
  (if (program-done? prog)
      #f
      (~>> prog
           (update-program (charterm-read-key))
           run-loop!)))

(define (program-done? prog)
  (~>> prog
       program-model
       ((program-exit? prog))))
  
(define (update-program key prog)
  (struct-copy program prog
               (model (~>> (program-model prog)
                           ((program-update prog) key)))))

(define (draw-program! prog)
  (charterm-clear-screen)
  (let-values [((w h) (charterm-screen-size))]
    (~>> (~> prog program-model)
         ((~> prog program-view) w h)
         (draw-lines! w h))))

(define (draw-lines! w h lines)
  ;; FIXME: build-list ?
  (~>> (iota h)
       (map (λ (i)
              (if (< i (length lines))
                  (list-ref lines i)
                  (make-string w #\space))))
       (for-each (λ (line) (draw-line! w line)))))

(define (draw-line! width line)
  (~>> (λ (i)
         (if (< i (string-length line))
             (string-ref line i)
             #\space))
       (build-string width)
       charterm-display))