WITH target_programs AS (
    SELECT id FROM programs
    WHERE id IN (
                 '2ZzfzVZMBaxfN8ZnraCZnP4JT28SvXH9ABjnwQgx36rn',
                 'J7Un1WWMe99WFQmR3gFeSvgbHGFvYiBaizkEXiTaaCiF',
                 '6jrnNGnETEnKzzQcPZrgoYVx3wGig6dKNdGgbTCtxx5o'
        )
),

     target_transactions AS (
         SELECT pt.transaction_signature AS signature
         FROM program_transactions pt
                  JOIN target_programs p ON pt.program_id = p.id
     )

DELETE FROM transactions
WHERE signature IN (SELECT signature FROM target_transactions);